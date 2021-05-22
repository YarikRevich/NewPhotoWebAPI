package models

type GETResponseEqualAlbumModel struct {
	Result struct {
		Name   string `json:"name"`
		Photos []struct {
			Photo     string `json:"photo"`
			Thumbnail string `json:"thumbnail"`
			Extension string `json:"extension"`
		} `json:"photos"`
		Videos []struct {
			Video     string `json:"video"`
			Extension string `json:"extension"`
		} `json:"videos"`
	} `json:"result"`
	Service struct {
		Ok      bool   `json:"ok"`
		Message string `json:"message"`
	} `json:"service"`
}

type POSTResponseEqualAlbumModel struct {
	Service struct {
		Ok      bool   `json:"ok"`
		Message string `json:"message"`
	} `json:"service"`
}

type PUTRequestEqualAlbumModel struct {
	Data struct {
		Name   string
		Photos []struct {
			File      string  `json:"file"`
			Name      string  `json:"name"`
			Size      float64 `json:"size"`
			Extension string  `json:"extension"`
		} `json:"photos"`
		Videos []struct {
			File      string  `json:"file"`
			Name      string  `json:"name"`
			Size      float64 `json:"size"`
			Extension string  `json:"extension"`
		} `json:"videos"`
	} `json:"data"`
}

type PUTResponseEqualAlbumModel struct {
	Service struct {
		Ok      bool   `json:"ok"`
		Message string `json:"message"`
	} `json:"service"`
}

type DELETERequestEqualAlbumModel struct {
	Data struct {
		Name      string  `json:"name"`
		Photos []string `json:"photos"`
		Videos []string `json:"videos"`
	} `json:"data"`
}

type DELETEResponseEqualAlbumModel struct {
	Service struct {
		Ok      bool   `json:"ok"`
		Message string `json:"message"`
	} `json:"service"`
}
