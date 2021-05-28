package models

type GETResponseEqualAlbumModel struct {
	Result struct {
		Name   string `json:"name"`
		Photos []struct {
			Photo     []byte `json:"photo"`
			Thumbnail []byte `json:"thumbnail"`
			Extension string `json:"extension"`
		} `json:"photos"`
		Videos []struct {
			Video     []byte `json:"video"`
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
	Service struct {
		Ok      bool   `json:"ok"`
		Message string `json:"message"`
	} `json:"service"`
}

type DELETERequestEqualAlbumModel struct {
	Data struct {
		Name      string  `json:"name"`
		Photos [][]byte `json:"photos"`
		Videos [][]byte `json:"videos"`
	} `json:"data"`
}

type DELETEResponseEqualAlbumModel struct {
	Service struct {
		Ok      bool   `json:"ok"`
		Message string `json:"message"`
	} `json:"service"`
}
