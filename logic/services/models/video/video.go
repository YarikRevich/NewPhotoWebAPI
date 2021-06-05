package video

import "NewPhotoWeb/logic/services/models/service"

type POSTRequestVideoModel struct {
	Data []struct {
		File      []byte  `json:"file"`
		Name      string  `json:"name"`
		Size      float64 `json:"size"`
		Extension string  `json:"extension"`
	} `json:"Data"`
}

type POSTResponseVideoModel struct {
	service.ServiceModel
}
