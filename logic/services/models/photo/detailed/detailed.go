package detailedphoto

import "NewPhotoWeb/logic/services/models/service"

type GETRequestDetailedPhotoModel struct {
	Data struct {
		Thumbnail []byte `json:"thumbnail"`
	} `json:"data"`
}

type GETResponseDetailedPhotoModel struct {
	Result struct {
		Photo []byte `json:"photo"`
	} `json:"result"`
	service.ServiceModel
}
