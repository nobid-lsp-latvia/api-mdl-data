// SPDX-License-Identifier: EUPL-1.2

package csdd

import (
	"azugo.io/core/validation"
	"github.com/spf13/viper"
)

type Configuration struct {
	CSDDUrl                string `mapstructure:"csdd_url"`
	CSDDUserName           string `mapstructure:"csdd_username"`
	CSDDChangePasswordDays string `mapstructure:"csdd_change_password_days"`
	CSDDSystemGUID         string `mapstructure:"csdd_system_guid"`
	CSDDSystemName         string `mapstructure:"csdd_system_name"`
	SkipVerify             bool   `mapstructure:"skip_verify"`
}

func (c *Configuration) Bind(prefix string, v *viper.Viper) {
	_ = v.BindEnv(prefix+".csdd_change_password_days", "CSDD_CHANGE_PASSWORD_DAYS")
	_ = v.BindEnv(prefix+".csdd_url", "CSDD_URL")
	_ = v.BindEnv(prefix+".csdd_username", "CSDD_USERNAME")
	_ = v.BindEnv(prefix+".csdd_change_password_days", "CSDD_CHANGE_PASSWORD_DAYS")
	_ = v.BindEnv(prefix+".csdd_system_guid", "CSDD_SYSTEM_GUID")
	_ = v.BindEnv(prefix+".csdd_system_name", "CSDD_SYSTEM_NAME")
	_ = v.BindEnv(prefix+".skip_verify", "CSDD_SKIP_TLS_VERIFY")
}

func (c *Configuration) Validate(valid *validation.Validate) error {
	return valid.Struct(c)
}
