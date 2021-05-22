package models

// type GETResponseAlbumModel struct {
// 	Result []struct {
// 		Photo     string   `json:"photo"`
// 		Thumbnail string   `json:"thumbnail"`
// 		Tags      []string `json:"tags"`
// 	} `json:"result"`
// 	Service struct {
// 		Message string `json:"message"`
// 		Ok      bool   `json:"ok"`
// 	} `json:"service"`
// }

type POSTRequestVideoModel struct {
	Data []struct {
		File      string  `json:"file"`
		Name      string  `json:"name"`
		Size      float64 `json:"size"`
		Extension string  `json:"extension"`
	} `json:"Data"`
}

type POSTResponseVideoModel struct {
	Service struct {
		Ok      bool   `json:"ok"`
		Message string `json:"message"`
	} `json:"service"`
}
