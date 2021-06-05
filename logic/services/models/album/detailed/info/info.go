package info

import "NewPhotoWeb/logic/services/models/service"

type GETRequestGetAlbumInfoModel struct {
	Data struct {
		Name string `json:"name"`
	} `json:"data"`
}

type GETResponseGetAlbumInfoModel struct {
	Result struct {
		MediaNum int64 `json:"media_num"`
	} `json:"result"`
	service.ServiceModel
}
