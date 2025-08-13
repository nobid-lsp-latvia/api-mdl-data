// SPDX-License-Identifier: EUPL-1.2

package responses

import "time"

type VaultGetTokenResponse struct {
	RequestID     string      `json:"request_id"`
	LeaseID       string      `json:"lease_id"`
	Renewable     bool        `json:"renewable"`
	LeaseDuration int         `json:"lease_duration"`
	Data          interface{} `json:"data"`
	WrapInfo      interface{} `json:"wrap_info"`
	Warnings      interface{} `json:"warnings"`
	Auth          struct {
		ClientToken   string   `json:"client_token"`
		Accessor      string   `json:"accessor"`
		Policies      []string `json:"policies"`
		TokenPolicies []string `json:"token_policies"`
		Metadata      struct {
			RoleName string `json:"role_name"`
		} `json:"metadata"`
		LeaseDuration  int         `json:"lease_duration"`
		Renewable      bool        `json:"renewable"`
		EntityID       string      `json:"entity_id"`
		TokenType      string      `json:"token_type"`
		Orphan         bool        `json:"orphan"`
		MfaRequirement interface{} `json:"mfa_requirement"`
		NumUses        int         `json:"num_uses"`
	} `json:"auth"`
	MountType string   `json:"mount_type"`
	Errors    []string `json:"errors"`
}

type VaultGetDataResponse struct {
	RequestID     string `json:"request_id"`
	LeaseID       string `json:"lease_id"`
	Renewable     bool   `json:"renewable"`
	LeaseDuration int    `json:"lease_duration"`
	Data          struct {
		Data struct {
			Password string `json:"edim-csdd-service-password"`
		} `json:"data"`
		Metadata struct {
			CreatedTime    time.Time   `json:"created_time"`
			CustomMetadata interface{} `json:"custom_metadata"`
			DeletionTime   string      `json:"deletion_time"`
			Destroyed      bool        `json:"destroyed"`
			Version        int         `json:"version"`
		} `json:"metadata"`
	} `json:"data"`
	WrapInfo  interface{} `json:"wrap_info"`
	Warnings  interface{} `json:"warnings"`
	Auth      interface{} `json:"auth"`
	MountType string      `json:"mount_type"`
	Errors    []string    `json:"errors"`
}

type VaultSaveDataPostResponse struct {
	RequestID     string `json:"request_id"`
	LeaseID       string `json:"lease_id"`
	Renewable     bool   `json:"renewable"`
	LeaseDuration int    `json:"lease_duration"`
	Data          struct {
		CreatedTime    time.Time   `json:"created_time"`
		CustomMetadata interface{} `json:"custom_metadata"`
		DeletionTime   string      `json:"deletion_time"`
		Destroyed      bool        `json:"destroyed"`
		Version        int         `json:"version"`
	} `json:"data"`
	WrapInfo  interface{} `json:"wrap_info"`
	Warnings  interface{} `json:"warnings"`
	Auth      interface{} `json:"auth"`
	MountType string      `json:"mount_type"`
	Errors    []string    `json:"errors"`
}

type VaultPostData struct {
	Data struct {
		Password string `json:"edim-csdd-service-password"`
	} `json:"data"`
}
