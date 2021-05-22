package models

type GETResponseAvatarModel struct {
	Result struct {
		Avatar string `json:"avatar"`
	} `json:"result"`
	Service struct {
		Ok bool `json:"ok"`
		Message string `json:"message"`
	} `json:"service"`
}

type POSTRequestAvatarModel struct {
	Data struct {
		Avatar string `json:"avatar"`
	} `json:"data"`
}

type POSTResponseAvatarModel struct {
	Service struct {
		Ok bool `json:"ok"`
		Message string `json:"message"`
	} `json:"service"`
}
