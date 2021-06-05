package account

import "NewPhotoWeb/logic/services/models/service"

type GETResponseAccountModel struct {
	Result struct {
		Firstname  string  `json:"firstname"`
		Secondname string  `json:"secondname"`
		Storage    float64 `json:"storage"`
	} `json:"result"`
	service.ServiceModel
}

type DELETEResponseAccountModel struct {
	service.ServiceModel
}
