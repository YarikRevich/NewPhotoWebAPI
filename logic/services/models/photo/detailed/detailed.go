package detailedphoto

import "NewPhotoWeb/logic/services/models/service"

type POSTRequestDetailedPhotoModel struct {
	Data struct {
		Thumbnail []byte `json:"thumbnail"`
	} `json:"data"`
}

type POSTResponseDetailedPhotoModel struct {
	Result struct {
		Media []byte `json:"media"`
	} `json:"result"`
	service.ServiceModel
}
