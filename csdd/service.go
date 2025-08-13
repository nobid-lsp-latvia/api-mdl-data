// SPDX-License-Identifier: EUPL-1.2

package csdd

import (
	"git.zzdats.lv/edim/api-mdl/routes/responses"
	"git.zzdats.lv/edim/api-mdl/vault"

	"azugo.io/azugo"
	"azugo.io/core"
)

type Service interface {
	GetCSDDData(ctx *azugo.Context, code string) (*responses.GetDataResponse, error)
}

func New(app *core.App, config *Configuration, vault vault.Service) (Service, error) {
	return newCsddService(app, config, vault)
}
