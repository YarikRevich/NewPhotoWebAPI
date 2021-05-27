package info

import (
	"context"
	"encoding/json"
	"net/http"

	. "NewPhotoWeb/config"
	"NewPhotoWeb/logic/proto"
	infodetailedalbummodel "NewPhotoWeb/logic/services/models/album/detailed/info"
	errormodel "NewPhotoWeb/logic/services/models/error"
)

type IInfoDetailedAlbumPage interface {
	GetHandler() http.Handler
}

type infodetailedalbum struct{}

func (a *infodetailedalbum) GetHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		values, ok := r.URL.Query()["name"]
		if !ok {
			Logger.Fatalln("Album name is empty!")
		}

		errResp := new(errormodel.ERRORAuthModel)
		errResp.Service.Error = errormodel.AUTH_ERROR
		at, err := r.Cookie("at")
		if err != nil {
			if err := json.NewEncoder(w).Encode(errResp); err != nil {
				Logger.Fatalln(err)
			}
			return
		}
		lt, err := r.Cookie("lt")
		if err != nil {
			if err := json.NewEncoder(w).Encode(errResp); err != nil {
				Logger.Fatalln(err)
			}
			return
		}

		grpcResp, err := NPC.GetAlbumInfo(
			context.Background(),
			&proto.GetAlbumInfoRequest{
				AccessToken: at.Value,
				LoginToken:  lt.Value,
				Album:       values[0],
			},
		)
		if err != nil {
			Logger.ClientError()
		}

		var resp infodetailedalbummodel.GETResponseGetAlbumInfoModel
		resp.Result.MediaNum = grpcResp.GetMediaNum()
		resp.Service.Ok = grpcResp.GetOk()

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			Logger.Fatalln(err)
		}
	})
}

func NewInfoDetailedAlbumPageHandler() IInfoDetailedAlbumPage {
	return new(infodetailedalbum)
}
