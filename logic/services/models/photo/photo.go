package photo

import "NewPhotoWeb/logic/services/models/service"

type GETResponsePhotoModel struct {
	Result []struct {
		Thumbnail []byte   `json:"thumbnail"`
		Tags      []string `json:"tags"`
	} `json:"result"`
	service.ServiceModel
}

type POSTRequestPhotoModel struct {
	Data []struct {
		File      []byte `json:"file"`
		Name      string  `json:"name"`
		Size      float64 `json:"size"`
		Extension string  `json:"extension"`
	} `json:"data"`
}

type POSTResponsePhotoModel struct {
	service.ServiceModel
}
