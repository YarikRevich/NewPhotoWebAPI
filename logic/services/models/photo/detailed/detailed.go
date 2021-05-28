package detailed

type GETRequestDetailedPhotoModel struct {
	Data struct {
		Thumbnail []byte `json:"thumbnail"`
	} `json:"data"`
}

type GETResponseDetailedPhotoModel struct {
	Result struct {
		Photo []byte `json:"photo"`
	} `json:"result"`
	Service struct {
		Ok      bool   `json:"ok"`
		Message string `json:"message"`
	} `json:"service"`
}
