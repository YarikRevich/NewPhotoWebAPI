package info

import (
	"context"
	"encoding/json"
	"net/http"


	"NewPhotoWeb/log"
	"NewPhotoWeb/logic/client"
	"NewPhotoWeb/logic/proto"
	infodetailedalbummodel "NewPhotoWeb/logic/services/models/album/detailed/info"
)

type IInfoDetailedAlbumPage interface {
	GetHandler() http.Handler
}

type infodetailedalbum struct{}

func (a *infodetailedalbum) GetHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		values, ok := r.URL.Query()["name"]
		if !ok {
			log.Logger.Fatalln("Album name is empty!")
		}

		at, _ := r.Cookie("at")
		lt, _ := r.Cookie("lt")

		grpcResp, err := client.NewPhotoClient.GetAlbumInfo(
			context.Background(),
			&proto.GetAlbumInfoRequest{
				AccessToken: at.Value,
				LoginToken:  lt.Value,
				Album:       values[0],
			},
		)
		if err != nil {
			log.Logger.ClientError(); client.Restart()
		}

		var resp infodetailedalbummodel.GETResponseGetAlbumInfoModel
		resp.Result.MediaNum = grpcResp.GetMediaNum()
		resp.Service.Ok = grpcResp.GetOk()

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Logger.Fatalln(err)
		}
	})
}

func NewInfoDetailedAlbumPageHandler() IInfoDetailedAlbumPage {
	return new(infodetailedalbum)
}
