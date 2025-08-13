// SPDX-License-Identifier: EUPL-1.2

package responses

import (
	"git.zzdats.lv/edim/api-mdl/utils"
)

// CategoryRestriction defines the driving privilege category restriction.
type CategoryRestriction struct {
	// Sign as per ISO/IEC 18013-2 Annex A
	Sign string `json:"sign"`
	// Value as per ISO/IEC 18013-2 Annex A
	Value string `json:"value"`
}

// DrivingPrivilege defines driving license categories.
type DrivingPrivilege struct {
	// Driver category code
	VehicleCategoryCode string `json:"vehicle_category_code"`
	// Starting date of validity of the driver category
	IssueDate utils.Date `json:"issue_date"`
	// Driver category expiry date
	ExpiryDate utils.Date `json:"expiry_date"`
	// Category restrictions
	Code []CategoryRestriction `json:"code"`
}

// MDLResponse defines the response structure for the CSDD data.
type MDLResponse struct {
	// PersonalAdministrativeNumber represents driver's administrative number
	PersonalAdministrativeNumber string `json:"personal_administrative_number"`
	// DocumentNumber represents document certificate number
	DocumentNumber string `json:"document_number"`
	// BirthData represents driver's date of birth
	BirthDate utils.Date `json:"birth_date"`
	// GivenName represents driver's name
	GivenName string `json:"given_name"`
	// FamilyName represents driver's surname
	FamilyName string `json:"family_name"`
	// IssueDate represents driver's license start date of validity
	IssueDate utils.Date `json:"issue_date"`
	// ExpireDate reprsents driver's license expiration date
	ExpiryDate utils.Date `json:"expiry_date"`
	// IssuingCountry represents driving license issuing country code (ISO 3166-1 alpha-2)
	IssuingCountry string `json:"issuing_country"`
	// IssuingAuthority represents issuing authority
	IssuingAuthority string `json:"issuing_authority"`
	// DrivingPrivileges represents driving license categories
	DrivingPrivileges []DrivingPrivilege `json:"driving_privileges"`
	// UnDistinguishingSign represents distinguishing sign of the issuing country according to ISO/IEC 18013-1:2018
	UnDistinguishingSign string `json:"un_distinguishing_sign"`
	// Portrait represent photo of the driver of the vehicle
	Portrait string `json:"portrait"`
}

type ChangePasswordResponse struct {
	Errors []ErrorResponse `json:"errors"`
}

type LoginResponse struct {
	Rowset []struct {
		SessionID string           `json:"sessionid"`
		PM        int              `json:"pm"`
		Fpda      bool             `json:"fpda"`
		Amats     string           `json:"amats"`
		KiestID   string           `json:"kiest_id"`
		KiestNos  string           `json:"kiest_nos"`
		KtKodi    []KtKodiResponse `json:"kt_kodi"`
		Kiest     []KiestResponse  `json:"kiest"`
	} `json:"rowset"`
	Errors []ErrorResponse `json:"errors"`
}

type KtKodiResponse struct {
	Kods string `json:"kods"`
	Ties bool   `json:"ties"`
}

type KiestResponse struct {
	ID  string `json:"id"`
	Nos string `json:"nos"`
}

type GetDataResponse struct {
	Rowset []*MDLResponse   `json:"rowset"`
	Errors []*ErrorResponse `json:"errors"`
}

type ErrorResponse struct {
	ClientMessageCode string `json:"clientMessageCode"`
	ClientMessage     string `json:"clientMessage"`
	Code              string `json:"code"`
}
