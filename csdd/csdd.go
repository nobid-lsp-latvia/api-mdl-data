// SPDX-License-Identifier: EUPL-1.2

package csdd

import (
	"crypto/rand"
	"errors"
	"math/big"
	"strconv"
	"sync"
	"time"

	"git.zzdats.lv/edim/api-mdl/routes/responses"
	"git.zzdats.lv/edim/api-mdl/vault"

	"azugo.io/azugo"
	"azugo.io/core"
	"azugo.io/core/http"
	"go.uber.org/zap"
)

const (
	charsNumbers = "0123456789"
	charsUppers  = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	charsLowers  = "abcdefghijklmnopqrstuvwxyz"
	charsSpecial = "~!@-#$+?"
)

type csddService struct {
	app    *core.App
	config *Configuration
	vault  vault.Service

	tokenMu sync.Mutex
}

func newCsddService(app *core.App, config *Configuration, vault vault.Service) (Service, error) {
	s := &csddService{
		app:    app,
		config: config,
		vault:  vault,
	}

	return s, nil
}

func (s *csddService) GetCSDDData(ctx *azugo.Context, code string) (*responses.GetDataResponse, error) {
	token, err := s.Login(ctx)
	if err != nil {
		return nil, err
	}
	defer s.Logout(ctx, token)

	response, err := s.GetData(ctx, token, code)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *csddService) Login(ctx *azugo.Context) (string, error) {
	s.tokenMu.Lock()
	defer s.tokenMu.Unlock()

	// get auth data from vault
	vaultdata, err := s.vault.GetCSDDAuthData(ctx, 0)
	if err != nil {
		return "", err
	}

	response, err := s.CallLogin(ctx, vaultdata)
	if err != nil {
		ctx.Log().Error("Finish get csdd sessionID with error", zap.Error(err))

		return "", err
	}

	// if is error in response, then check if error code is F-00011
	if len(response.Errors) > 0 {
		// if error code= F-00011, get one prior password,
		if response.Errors[0].ClientMessageCode == "F-00011" {
			vaultdataPriorVersion, err := s.vault.GetCSDDAuthData(ctx, vaultdata.Data.Metadata.Version)
			if err != nil {
				return "", err
			}

			// try login with prior password
			response, err = s.CallLogin(ctx, vaultdataPriorVersion)
			// if can not login with prior password, then panic
			if err != nil {
				return "", err
			}

			// if ok login with prior password, then save prior correct password to vault
			if _, err := s.vault.ChangeVaultData(ctx, vaultdataPriorVersion.Data.Data.Password); err != nil {
				return "", err
			}

			// create and save new password to vault and csdd
			s.ChangePassword(ctx, vaultdataPriorVersion, response.Rowset[0].SessionID)

			// logout from incorrect session
			// nevaig -> s.Logout(ctx, response.Rowset[0].SessionID)

			// get leatest vault data
			vaultdata, err = s.vault.GetCSDDAuthData(ctx, 0)
			if err != nil {
				return "", err
			}
			// call login with new password
			response, err = s.CallLogin(ctx, vaultdata)
			if err != nil {
				return "", err
			}
		} else {
			ctx.Log().Error("Error login to CSDD: " + response.Errors[0].ClientMessageCode + ": " + response.Errors[0].ClientMessage)

			return "", errors.New(response.Errors[0].ClientMessageCode + ": " + response.Errors[0].ClientMessage)
		}
	}

	// if PM == 1 vai PM == 2, need new password
	// or if password is older than "s.config.ChangePasswordDays" days
	// but we can call data by this session
	Days, _ := strconv.Atoi(s.config.CSDDChangePasswordDays)
	DurationDays := time.Duration(24*Days) * time.Hour

	if response.Rowset[0].PM == 1 ||
		response.Rowset[0].PM == 2 ||
		// ja parole ir vecāka par "s.config.ChangePasswordDays" dienām
		vaultdata.Data.Metadata.CreatedTime.Add(DurationDays).Before(time.Now()) {
		s.ChangePassword(ctx, vaultdata, response.Rowset[0].SessionID)
	}

	return response.Rowset[0].SessionID, nil
}

func (s *csddService) CallLogin(ctx *azugo.Context, vaultdata *responses.VaultGetDataResponse) (*responses.LoginResponse, error) {
	response := &responses.LoginResponse{}

	client := ctx.HTTPClient()

	if s.config.SkipVerify {
		client = ctx.HTTPClient().WithOptions(&http.TLSConfig{InsecureSkipVerify: true})
	}

	err := client.PostJSON(
		s.config.CSDDUrl,
		struct {
			SystemGUID  string `json:"SystemGUID"`
			SystemName  string `json:"SystemName"`
			ServiceName string `json:"ServiceName"`
			Params      struct {
				LietVards string `json:"liet_vards"`
				Parole    string `json:"parole"`
			} `json:"Params"`
		}{
			SystemGUID:  s.config.CSDDSystemGUID,
			SystemName:  s.config.CSDDSystemName,
			ServiceName: "Chk_web_Gliet",
			Params: struct {
				LietVards string `json:"liet_vards"`
				Parole    string `json:"parole"`
			}{
				LietVards: s.config.CSDDUserName,
				Parole:    vaultdata.Data.Data.Password,
			},
		},
		response,
	)
	if err != nil {
		ctx.Log().Error("Finish get csdd sessionID with error", zap.Error(err))

		return nil, err
	}

	return response, nil
}

func (s *csddService) Logout(ctx *azugo.Context, token string) {
	client := ctx.HTTPClient()

	if s.config.SkipVerify {
		client = ctx.HTTPClient().WithOptions(&http.TLSConfig{InsecureSkipVerify: true})
	}

	ctx.Log().Debug("===> start csdd logout")

	if err := client.PostJSON(
		s.config.CSDDUrl,
		struct {
			SystemGUID  string `json:"SystemGUID"`
			SystemName  string `json:"SystemName"`
			ServiceName string `json:"Del_Fses"`
			SessionID   string `json:"SessionID"`
		}{
			SystemGUID:  s.config.CSDDSystemGUID,
			SystemName:  s.config.CSDDSystemName,
			ServiceName: "Del_Fses",
			SessionID:   token,
		},
		nil,
	); err != nil {
		ctx.Log().Error("Finish csdd logout with error", zap.Error(err))
	}

	ctx.Log().Debug("===> finish csdd logout")
}

func (s *csddService) GetData(ctx *azugo.Context, token string, code string) (*responses.GetDataResponse, error) {
	response := &responses.GetDataResponse{}
	client := ctx.HTTPClient()

	if s.config.SkipVerify {
		client = ctx.HTTPClient().WithOptions(&http.TLSConfig{InsecureSkipVerify: true})
	}

	ctx.Log().Debug("===> start get csdd data")

	err := client.PostJSON(
		s.config.CSDDUrl,
		struct {
			SystemGUID  string `json:"SystemGUID"`
			SystemName  string `json:"SystemName"`
			ServiceName string `json:"ServiceName"`
			SessionID   string `json:"SessionID"`
			Params      struct {
				Pk   string `json:"pk"`
				Num  string `json:"num"`
				Foto bool   `json:"foto"`
			} `json:"Params"`
		}{
			SystemGUID:  s.config.CSDDSystemGUID,
			SystemName:  s.config.CSDDSystemName,
			ServiceName: "Qry_va",
			SessionID:   token,
			Params: struct {
				Pk   string `json:"pk"`
				Num  string `json:"num"`
				Foto bool   `json:"foto"`
			}{
				Pk:   code,
				Num:  "",
				Foto: true,
			},
		},
		response,
	)
	if err != nil {
		ctx.Log().Error("Finish get csdd data with error", zap.Error(err))

		return nil, err
	}

	return response, nil
}

func (s *csddService) ChangePassword(ctx *azugo.Context, indata *responses.VaultGetDataResponse, sessionID string) {
	// vispirms saglabājam jauno paroli Vault
	oldPsw := indata.Data.Data.Password
	newPsw := generateNewPassword()

	res, err := s.vault.ChangeVaultData(ctx, newPsw)
	// if Vault success, call CSDD change password
	if err == nil && len(res.Errors) == 0 {
		result, err := s.CallChangePassword(ctx, indata, sessionID, newPsw)
		// when succes, response ir empty

		// if error when change password in CSDD
		if err != nil || len(result.Errors) > 0 {
			ctx.Log().Error("Error changing password in CSDD", zap.Error(err)) // change back to old password in vault
			_, _ = s.vault.ChangeVaultData(ctx, oldPsw)                        // if error in Vault, then next login is with error "F-00011"
		}
	}
}

func (s *csddService) CallChangePassword(ctx *azugo.Context, indata *responses.VaultGetDataResponse, sessionID string, newPsw string) (*responses.ChangePasswordResponse, error) {
	result := &responses.ChangePasswordResponse{}
	client := ctx.HTTPClient()

	err := client.PostJSON(
		s.config.CSDDUrl,
		struct {
			SystemGUID  string `json:"SystemGUID"`
			SystemName  string `json:"SystemName"`
			ServiceName string `json:"ServiceName"`
			SessionID   string `json:"SessionID"`
			Params      struct {
				IeprParole string `json:"iepr_parole"`
				Parole     string `json:"parole"`
			} `json:"Params"`
		}{
			SystemGUID:  s.config.CSDDSystemGUID,
			SystemName:  s.config.CSDDSystemName,
			ServiceName: "Upd_web_parole",
			SessionID:   sessionID,
			Params: struct {
				IeprParole string `json:"iepr_parole"`
				Parole     string `json:"parole"`
			}{
				IeprParole: indata.Data.Data.Password,
				Parole:     newPsw,
			},
		},
		result,
	)
	if err != nil {
		ctx.Log().Error("Finish change password with error", zap.Error(err))

		return result, err
	}

	return result, nil
}

/**
 * generateNewPassword is a function that generates a new password.
 */
func generateNewPassword() string {
	b := make([]byte, 16)

	r := map[int]int{0: 2}
	for i := range r {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charsNumbers))))
		b[i] = charsNumbers[num.Int64()]
	}

	for i := 2; i < 5; i++ {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charsUppers))))
		b[i] = charsUppers[num.Int64()]
	}

	for i := 5; i < 14; i++ {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charsLowers))))
		b[i] = charsLowers[num.Int64()]
	}

	for i := 14; i < 16; i++ {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charsSpecial))))
		b[i] = charsSpecial[num.Int64()]
	}

	for i := len(b) - 1; i > 0; i-- {
		jBig, _ := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		j := int(jBig.Int64())
		b[i], b[j] = b[j], b[i]
	}

	return string(b)
}
