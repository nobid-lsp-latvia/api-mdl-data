// SPDX-License-Identifier: EUPL-1.2

package main

import (
	"fmt"
	"os"

	app "git.zzdats.lv/edim/api-mdl"

	"github.com/spf13/cobra"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

type HealthCheck struct {
	Path    string
	Version string
}

func (h *HealthCheck) RunHealthCheck(cmd *cobra.Command) error {
	a, err := app.New(cmd, h.Version)
	if err != nil {
		a.Log().Error("failed to load configuration", zap.Error(err))
		os.Exit(1)

		return nil
	}

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.Header.SetMethod(fasthttp.MethodGet)
	req.SetRequestURI(fmt.Sprintf("http://localhost:%d%s", a.Config().Server.HTTP.Port, h.Path))

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	client := &fasthttp.Client{}

	return client.Do(req, resp)
}
