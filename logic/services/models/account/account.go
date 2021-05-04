package models

type GETResponseAccountModel struct {
	Result struct {
		Firstname  string  `json:"firstname"`
		Secondname string  `json:"secondname"`
		Storage    float64 `json:"storage"`
	} `json:"result"`
	Service struct {
		Ok      bool   `json:"ok"`
		Message string `json:"message"`
	} `json:"service"`
}
