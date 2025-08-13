// SPDX-License-Identifier: EUPL-1.2

package routes

import (
	app "git.zzdats.lv/edim/api-mdl"
	"git.zzdats.lv/edim/api-mdl/openapi"

	"github.com/nobid-lsp-latvia/go-idauth"
	oa "github.com/nobid-lsp-latvia/go-openapi"
)

type router struct {
	*app.App
	openapi *oa.OpenAPI
}

func Init(a *app.App) error {
	r := &router{
		App: a,
	}
	r.openapi = oa.NewDefaultOpenAPIHandler(openapi.OpenAPIDefinition, a.App)

	a.Get("/healthz", r.healthz)

	v1 := a.Group("/1.0")
	{
		v1.Use(idauth.Authentication(a.App, a.Config().IDAuth))

		v1.Get("/mdl", idauth.UserHasScope("citizen", r.mdl))
	}

	return nil
}
