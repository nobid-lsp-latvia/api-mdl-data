// SPDX-License-Identifier: EUPL-1.2

package vault

import (
	"azugo.io/core/config"
	"azugo.io/core/validation"
	"github.com/spf13/viper"
)

// Configuration represents the configuration for the vault service.
type Configuration struct {
	LoginURL string `mapstructure:"url_login"`
	DataURL  string `mapstructure:"url_data"`
	RoleID   string `mapstructure:"role_id"`
	SecretID string `mapstructure:"secret_id"`
}

func (c *Configuration) Bind(prefix string, v *viper.Viper) {
	key, _ := config.LoadRemoteSecret("VAULT_SECRET_ID")
	v.SetDefault(prefix+".secret_id", key)

	_ = v.BindEnv(prefix+".url_login", "VAULT_LOGIN_URL")
	_ = v.BindEnv(prefix+".url_data", "VAULT_DATA_URL")
	_ = v.BindEnv(prefix+".role_id", "VAULT_ROLE_ID")
	_ = v.BindEnv(prefix+".secret_id", "VAULT_SECRET_ID")
}

// Validate vault configuration section.
func (c *Configuration) Validate(valid *validation.Validate) error {
	return valid.Struct(c)
}
