package signup

import "NewPhotoWeb/logic/services/models/service"

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
	service.ServiceModel
}
