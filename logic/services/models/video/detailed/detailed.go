package detailedvideo

import "NewPhotoWeb/logic/services/models/service"

type GETResponseDetailedVideoModel struct {
	Result struct {
		Photo []byte `json:"photo"`
	} `json:"result"`
	service.ServiceModel
}
