package models

type GETResponseAlbumModel struct {
	Result []struct {
		Name                 string `json:"name"`
		LatestPhoto          string `json:"latestphoto"`
		LatestPhotoThumbnail string `json:"latestphotothumbnail"`
	} `json:"result"`
	Service struct {
		Ok      bool   `json:"ok"`
		Message string `json:"message"`
	} `json:"service"`
}

type POSTRequestAlbumModel struct {
	Data struct {
		Name string `json:"name"`
	} `json:"data"`
}

type POSTResponseAlbumModel struct {
	Service struct {
		Ok      bool   `json:"ok"`
		Message string `json:"message"`
	} `json:"service"`
}

type DELETERequestAlbumModel struct {
	Data struct {
		Name string `json:"name"`
	} `json:"data"`
}

type DELETEResponseAlbumModel struct {
	Service struct {
		Ok      bool   `json:"ok"`
		Message string `json:"message"`
	} `json:"service"`
}

type PUTRequestAlbumModel struct {
	Result struct {
		Name string `json:"name"`
		Data []struct {
			File      string  `json:"file"`
			Size      float64 `json:"size"`
			Extension string  `json:"extension"`
		} `json:"data"`
	} `json:"result"`
}

type PUTResponseAlbumModel struct {
	Service struct {
		Ok      bool   `json:"ok"`
		Message string `json:"message"`
	} `json:"service"`
}
