package models

type GETResponseSignOutModel struct {
	Service struct {
		Ok bool `json:"ok"`
	} `json:"service"`
}
