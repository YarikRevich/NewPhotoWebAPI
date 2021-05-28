package models

type POSTRequestVideoModel struct {
	Data []struct {
		File      []byte  `json:"file"`
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
