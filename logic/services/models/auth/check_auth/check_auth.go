package models

type GETResponseCheckAuthModel struct {
	Service struct {
		Ok      bool   `json:"ok"`
	} `json:"service"`
}
