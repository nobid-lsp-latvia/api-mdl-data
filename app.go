// SPDX-License-Identifier: EUPL-1.2

package mdl

import (
	"git.zzdats.lv/edim/api-mdl/csdd"
	"git.zzdats.lv/edim/api-mdl/vault"

	"azugo.io/azugo"
	"azugo.io/azugo/server"
	"github.com/spf13/cobra"
)

// App is the application instance.
type App struct {
	*azugo.App

	config *Configuration
	vault  vault.Service
	csdd   csdd.Service
}

// New returns a new application instance.
func New(cmd *cobra.Command, version string) (*App, error) {
	config := NewConfiguration()

	a, err := server.New(cmd, server.Options{
		AppName:       "API-MDL",
		AppVer:        version,
		Configuration: config,
	})
	if err != nil {
		return nil, err
	}

	instance := &App{
		App:    a,
		config: config,
	}

	err = instance.InitServices()
	if err != nil {
		return nil, err
	}

	return instance, nil
}

func (a *App) InitServices() error {
	var err error

	a.vault, err = vault.New(a.App.App, a.config.Vault)
	if err != nil {
		return err
	}

	a.csdd, err = csdd.New(a.App.App, a.config.CSDD, a.vault)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) VaultService() vault.Service {
	return a.vault
}

func (a *App) CsddService() csdd.Service {
	return a.csdd
}

// Config returns application configuration.
//
// Panics if configuration is not loaded.
func (a *App) Config() *Configuration {
	if a.config == nil || !a.config.Ready() {
		panic("configuration is not loaded")
	}

	return a.config
}
