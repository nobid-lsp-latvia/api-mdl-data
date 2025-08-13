// SPDX-License-Identifier: EUPL-1.2

package mdl

import (
	"time"

	"git.zzdats.lv/edim/api-mdl/csdd"
	"git.zzdats.lv/edim/api-mdl/vault"

	"azugo.io/azugo/config"
	"azugo.io/core/validation"
	"github.com/nobid-lsp-latvia/go-idauth"
	"github.com/spf13/viper"
)

// Configuration represents the configuration for the application.
type Configuration struct {
	*config.Configuration `mapstructure:",squash"`

	Vault  *vault.Configuration  `mapstructure:"vault"`
	CSDD   *csdd.Configuration   `mapstructure:"csdd"`
	IDAuth *idauth.Configuration `mapstructure:"idauth"`
}

// NewConfiguration returns a new configuration.
func NewConfiguration() *Configuration {
	return &Configuration{
		Configuration: config.New(),
	}
}

// Core returns the core configuration.
func (c *Configuration) ServerCore() *config.Configuration {
	return c.Configuration
}

// Bind configuration to viper.
func (c *Configuration) Bind(_ string, v *viper.Viper) {
	c.Configuration.Bind("", v)

	c.Vault = config.Bind(c.Vault, "vault", v)
	c.CSDD = config.Bind(c.CSDD, "csdd", v)
	c.IDAuth = config.Bind(c.IDAuth, "idauth", v)
}

// Validate application configuration.
func (c *Configuration) Validate(validate *validation.Validate) error {
	if err := c.Vault.Validate(validate); err != nil {
		return err
	}

	if err := c.CSDD.Validate(validate); err != nil {
		return err
	}

	if err := c.IDAuth.Validate(validate); err != nil {
		return err
	}

	return nil
}

type ClientConfiguration struct {
	SessionTimeout   time.Duration `mapstructure:"session_timeout" validate:"gt=0,omitempty"`
	SessionCountdown time.Duration `mapstructure:"session_countdown"`
}

func (c *ClientConfiguration) Bind(prefix string, v *viper.Viper) {
	// vault un csdd timeout itkā ir 24H, ieliksim pagaidām 2h
	v.SetDefault(prefix+".session_timeout", 2*time.Hour)

	_ = v.BindEnv(prefix+".session_timeout", "SESSION_TIMEOUT")
}
