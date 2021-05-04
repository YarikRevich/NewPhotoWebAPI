package models

type GetRespEqualAlbumModel struct {
	Name   string `json:"name"`
	Result []struct {
		Photo string `json:"photo"`
		Thumbnail string `json:"thumbnail"`
	} `json:"result"`
	Service struct {
		Ok      bool   `json:"ok"`
		Message string `json:"message"`
	} `json:"service"`
}
