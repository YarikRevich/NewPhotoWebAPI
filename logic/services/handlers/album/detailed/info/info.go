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

		at := r.Header["X-At"]
		lt := r.Header["X-Lt"]

		grpcRespPhotos, err := client.NewPhotoClient.GetPhotosInAlbumNum(
			context.Background(),
			&proto.GetPhotosInAlbumNumRequest{
				AccessToken: at[0],
				LoginToken:  lt[0],
				Name:        values[0],
			},
		)
		if err != nil {
			log.Logger.ClientError()
			client.Restart()
		}

		grpcRespVideos, err := client.NewPhotoClient.GetVideosInAlbumNum(
			context.Background(),
			&proto.GetVideosInAlbumNumRequest{
				AccessToken: at[0],
				LoginToken:  lt[0],
				Name:        values[0],
			},
		)
		if err != nil {
			log.Logger.ClientError()
			client.Restart()
		}

		var resp infodetailedalbummodel.GETResponseGetAlbumInfoModel
		resp.Result.MediaNum = grpcRespPhotos.GetNum() + grpcRespVideos.GetNum()
		resp.Service.Ok = grpcRespPhotos.GetOk() && grpcRespVideos.GetOk()

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Logger.Fatalln(err)
		}
	})
}

func NewInfoDetailedAlbumPageHandler() IInfoDetailedAlbumPage {
	return new(infodetailedalbum)
}
