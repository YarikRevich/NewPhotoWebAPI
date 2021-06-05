package avatar

import "NewPhotoWeb/logic/services/models/service"

type GETResponseAvatarModel struct {
	Result struct {
		Avatar []byte `json:"avatar"`
	} `json:"result"`
	service.ServiceModel
}

type POSTRequestAvatarModel struct {
	Data struct {
		Avatar []byte `json:"avatar"`
	} `json:"data"`
}

type POSTResponseAvatarModel struct {
	service.ServiceModel
}
