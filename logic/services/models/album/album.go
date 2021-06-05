package albums

import "NewPhotoWeb/logic/services/models/service"

type GETResponseAlbumModel struct {
	Result []struct {
		Name                 string `json:"name"`
		LatestPhoto          []byte `json:"latestphoto"`
		LatestPhotoThumbnail []byte `json:"latestphotothumbnail"`
	} `json:"result"`
	service.ServiceModel
}

type POSTRequestAlbumModel struct {
	Data struct {
		Name string `json:"name"`
	} `json:"data"`
}

type POSTResponseAlbumModel struct {
	service.ServiceModel
}

type DELETERequestAlbumModel struct {
	Data struct {
		Name string `json:"name"`
	} `json:"data"`
}

type DELETEResponseAlbumModel struct {
	service.ServiceModel
}

type PUTRequestAlbumModel struct {
	Result struct {
		Name string `json:"name"`
		Data []struct {
			File      []byte  `json:"file"`
			Size      float64 `json:"size"`
			Extension string  `json:"extension"`
		} `json:"data"`
	} `json:"result"`
}

type PUTResponseAlbumModel struct {
	service.ServiceModel
}
