package video

import "NewPhotoWeb/logic/services/models/service"

type GETResponseVideoModel struct {
	Result []struct {
		Thumbnail []byte   `json:"thumbnail"`
		Tags      []string `json:"tags"`
	} `json:"result"`
	service.ServiceModel
}

type POSTRequestVideoModel struct {
	Data []struct {
		File      []uint8 `json:"file"`
		Name      string  `json:"name"`
		Size      float64 `json:"size"`
		Extension string  `json:"extension"`
	} `json:"Data"`
}

type POSTResponseVideoModel struct {
	service.ServiceModel
}
