// SPDX-License-Identifier: EUPL-1.2

package requests

type CSDDDataIn struct {
	SystemGUID  string `json:"SystemGUID"`
	SystemName  string `json:"SystemName"`
	ServiceName string `json:"ServiceName"`
	SessionID   string `json:"SessionID"`
	Params      Params `json:"Params"`
}

type Params struct {
	LietVards  string `json:"liet_vards"`
	IeprParole string `json:"iepr_parole"`
	Parole     string `json:"parole"`
	Pk         string `json:"pk"`
	Num        string `json:"num"`
	Foto       bool   `json:"foto"`
}
