package detailed

import (
	"context"
	"encoding/json"
	"net/http"

	"NewPhotoWeb/logic/proto"
	detailedvideomodel "NewPhotoWeb/logic/services/models/video/detailed"

	"NewPhotoWeb/log"
	"NewPhotoWeb/logic/client"
)

type IDetailedVideo interface {
	GetHandler() http.Handler
}

type detailedvideo struct{}

func (a *detailedvideo) GetHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t, ok := r.URL.Query()["thumbnail"]
		if !ok {
			log.Logger.Fatalln("Thumbnail is not passed")
		}

		at := r.Header["X-At"]
		lt := r.Header["X-Lt"]

		grpcResp, err := client.NewPhotoClient.GetFullMediaByThumbnail(
			context.Background(),
			&proto.GetFullMediaByThumbnailRequest{
				AccessToken: at[0],
				LoginToken:  lt[0],
				Thumbnail:   []byte(t[0]),
				MediaType:   proto.MediaType_Video,
			},
		)
		if err != nil {
			log.Logger.ClientError()
			client.Restart()
		}

		resp := new(detailedvideomodel.GETResponseDetailedVideoModel)
		resp.Service.Ok = true
		resp.Result.Photo = grpcResp.GetMedia()

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Logger.Fatalln(err)
		}
	})

}

func NewDetailedVideoHandler() IDetailedVideo {
	return new(detailedvideo)
}
