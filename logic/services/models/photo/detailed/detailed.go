package detailed

type GETRequestDetailedPhotoModel struct {
	Data struct {
		Photo string `json:"photo"`
	} `json:"data"`
}

type GETResponseDetailedPhotoModel struct {
	Result struct {
		Photo string `json:"photo"`
	} `json:"result"`
	Service struct {
		Ok      bool   `json:"ok"`
		Message string `json:"message"`
	} `json:"service"`
}