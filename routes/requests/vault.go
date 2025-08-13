// SPDX-License-Identifier: EUPL-1.2

package requests

type VaultGetTokenRequest struct {
	RoleID   string `json:"role_id"`
	SecretID string `json:"secret_id"`
}

type VaultSaveDataPostRequest struct {
	Data struct {
		UserName  string `json:"user_name"`
		Psw       string `json:"psw"`
		ValidDays int    `json:"valid_days"`
	} `json:"data"`
}
