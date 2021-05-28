package models

type GETResponsePhotoModel struct {
	Result []struct {
		Thumbnail []byte   `json:"thumbnail"`
		Tags      []string `json:"tags"`
	} `json:"result"`
	Service struct {
		Message string `json:"message"`
		Ok      bool   `json:"ok"`
	} `json:"service"`
}

type POSTRequestPhotoModel struct {
	Data []struct {
		File      []byte  `json:"file"`
		Name      string  `json:"name"`
		Size      float64 `json:"size"`
		Extension string  `json:"extension"`
	} `json:"Data"`
}

type POSTResponsePhotoModel struct {
	Service struct {
		Ok      bool   `json:"ok"`
		Message string `json:"message"`
	} `json:"service"`
}
