package models

type POSTRequestRegestrationModel struct {
	Data struct {
		Login      string `json:"login"`
		Firstname  string `json:"firstname"`
		Secondname string `json:"secondname"`
		Password1  string `json:"password1"`
		Password2  string `json:"password2"`
	} `json:"data"`
}

type POSTResponseRegestrationModel struct {
	Service struct {
		Ok bool `json:"ok"`
	} `json:"service"`
}
