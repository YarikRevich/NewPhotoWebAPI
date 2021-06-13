package detailedvideo

import "NewPhotoWeb/logic/services/models/service"

type POSTRequestDetailedVideoModel struct {
	Data struct {
		Thumbnail []byte `json:"thumbnail"`
	} `json:"data"`
}

type POSTResponseDetailedVideoModel struct {
	Result struct {
		Media []byte `json:"media"`
	} `json:"result"`
	service.ServiceModel
}