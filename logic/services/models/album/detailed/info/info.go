package info

type GETRequestGetAlbumInfoModel struct {
	Data struct {
		Name string `json:"name"`
	} `json:"data"`
}

type GETResponseGetAlbumInfoModel struct {
	Result struct {
		MediaNum int64 `json:"media_num"`
	} `json:"result"`
	Service struct {
		Ok      bool   `json:"ok"`
		Message string `json:"message"`
	} `json:"service"`
}
