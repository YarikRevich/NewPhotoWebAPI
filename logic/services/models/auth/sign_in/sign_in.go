package models

type GETRequestSignInModel struct {
	Data struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	} `json:"data"`
}

type GETResponseSignInModel struct {
	Service struct {
		Ok      bool   `json:"ok"`
		Message string `json:"message"`
	} `json:"service"`
}
