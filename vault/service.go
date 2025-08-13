// SPDX-License-Identifier: EUPL-1.2

package vault

import (
	"git.zzdats.lv/edim/api-mdl/routes/responses"

	"azugo.io/azugo"
	"azugo.io/core"
)

type Service interface {
	GetCSDDAuthData(ctx *azugo.Context, version int) (*responses.VaultGetDataResponse, error)
	ChangeVaultData(ctx *azugo.Context, newpsw string) (*responses.VaultSaveDataPostResponse, error)
}

func New(app *core.App, config *Configuration) (Service, error) {
	return newVaultService(app, config)
}
