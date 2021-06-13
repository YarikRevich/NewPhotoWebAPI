package equalalbum

import "NewPhotoWeb/logic/services/models/service"

type GETResponseEqualAlbumModel struct {
	Result struct {
		Name   string `json:"name"`
		Photos []struct {
			Thumbnail []byte `json:"thumbnail"`
			Extension string `json:"extension"`
		} `json:"photos"`
		Videos []struct {
			Video     []byte `json:"video"`
			Extension string `json:"extension"`
		} `json:"videos"`
	} `json:"result"`
	service.ServiceModel
}

type POSTResponseEqualAlbumModel struct {
	service.ServiceModel
}

type PUTRequestEqualAlbumModel struct {
	Data struct {
		Name   string
		Photos []struct {
			File      []byte  `json:"file"`
			Size      float64 `json:"size"`
			Extension string  `json:"extension"`
		} `json:"photos"`
		Videos []struct {
			File      []byte  `json:"file"`
			Size      float64 `json:"size"`
			Extension string  `json:"extension"`
		} `json:"videos"`
	} `json:"data"`
}

type PUTResponseEqualAlbumModel struct {
	service.ServiceModel
}

type DELETERequestEqualAlbumModel struct {
	Data struct {
		Name   string   `json:"name"`
		Photos [][]byte `json:"photos"`
		Videos [][]byte `json:"videos"`
	} `json:"data"`
}

type DELETEResponseEqualAlbumModel struct {
	service.ServiceModel
}
