package vault

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"git.zzdats.lv/edim/api-mdl/routes/responses"

	"azugo.io/azugo"
	"azugo.io/core"
	"azugo.io/core/cache"
	"azugo.io/core/http"
)

type vaultService struct {
	app     *core.App
	config  *Configuration
	cache   cache.Instance[string]
	tokenMu sync.RWMutex
}

func newVaultService(app *core.App, config *Configuration) (Service, error) {
	cache, err := cache.Create[string](cache.New(), "csdd-vault-token", cache.MemoryCache, cache.DefaultTTL(8*time.Hour))
	if err != nil {
		return nil, err
	}

	s := &vaultService{
		app:    app,
		config: config,
		cache:  cache,
	}

	return s, nil
}

func (s *vaultService) GetToken(ctx *azugo.Context) (string, error) {
	s.tokenMu.RLock()
	defer s.tokenMu.RUnlock()

	token, err := s.cache.Get(ctx, "secret_id")
	if err == nil && token != "" {
		return token, nil
	}

	s.tokenMu.RUnlock()
	s.tokenMu.Lock()

	response := &responses.VaultGetTokenResponse{}
	client := ctx.HTTPClient()

	ctx.Log().Debug("===> start get vault token")

	err = client.PostJSON(
		s.config.LoginURL, // "https://vault.zzdats.lv/v1/auth/lvrtc-edim/login",
		struct {
			RoleID   string `json:"role_id"`
			SecretID string `json:"secret_id"`
		}{
			RoleID:   s.config.RoleID,
			SecretID: s.config.SecretID,
		},
		response,
	)
	if err != nil {
		s.tokenMu.Unlock()
		s.tokenMu.RLock()

		return "", err
	}

	if len(response.Errors) > 0 {
		s.tokenMu.Unlock()
		s.tokenMu.RLock()

		return "", fmt.Errorf("vault error: %s", response.Errors[0])
	}

	if err := s.cache.Set(ctx, "secret_id", response.Auth.ClientToken); err != nil {
		return "", err
	}

	s.tokenMu.Unlock()
	s.tokenMu.RLock()

	ctx.Log().Debug("===> finish get vault token")

	return response.Auth.ClientToken, nil
}

func (s *vaultService) GetCSDDAuthData(ctx *azugo.Context, version int) (*responses.VaultGetDataResponse, error) {
	token, err := s.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	response, err := s.getVaultCSDDAuthData(ctx, token, version)

	return response, err
}

func (s *vaultService) getVaultCSDDAuthData(ctx *azugo.Context, token string, version int) (*responses.VaultGetDataResponse, error) {
	response := &responses.VaultGetDataResponse{}

	client := ctx.HTTPClient()
	link := s.config.DataURL

	if version > 0 {
		link = link + "?version=" + strconv.Itoa(version-1)
	}

	maxRetries := map[int]int{0: 3}

	var (
		lastErr error
		retry   int
	)

	for retry = range maxRetries {
		err := client.GetJSON(
			link,
			response,
			http.WithHeader("X-Vault-Token", token),
		)
		if err == nil {
			return response, nil
		}

		lastErr = err

		if !strings.Contains(err.Error(), "invalid token") {
			break
		}

		token, err = s.GetToken(ctx)
		if err != nil {
			continue
		}
	}

	return nil, fmt.Errorf("failed after %d retries: %w", retry, lastErr)
}

func (s *vaultService) ChangeVaultData(ctx *azugo.Context, newPsw string) (*responses.VaultSaveDataPostResponse, error) {
	result := &responses.VaultSaveDataPostResponse{}

	client := ctx.HTTPClient()

	token, err := s.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	postData := &responses.VaultPostData{}
	postData.Data.Password = newPsw

	err = client.PostJSON(
		s.config.DataURL,
		postData,
		result,
		http.WithHeader("X-Vault-Token", token),
	)
	if err != nil {
		return nil, err
	}

	return result, nil
}
