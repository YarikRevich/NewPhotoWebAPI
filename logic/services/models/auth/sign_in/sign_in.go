package signin

import "NewPhotoWeb/logic/services/models/service"

type GETRequestSignInModel struct {
	Data struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	} `json:"data"`
}

type GETResponseSignInModel struct {
	service.ServiceModel
}
