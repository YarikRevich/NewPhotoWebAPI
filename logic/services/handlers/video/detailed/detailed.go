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
	PostHandler() http.Handler
}

type detailedvideo struct{}

func (a *detailedvideo) PostHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var stat detailedvideomodel.POSTRequestDetailedVideoModel
		if err := json.NewDecoder(r.Body).Decode(&stat); err != nil {
			log.Logger.Fatalln(err)
		}

		at := r.Header["X-At"]
		lt := r.Header["X-Lt"]

		grpcResp, err := client.NewPhotoClient.GetFullMediaByThumbnail(
			context.Background(),
			&proto.GetFullMediaByThumbnailRequest{
				AccessToken: at[0],
				LoginToken:  lt[0],
				Thumbnail:   stat.Data.Thumbnail,
				MediaType:   proto.MediaType_Video,
			},
		)
		if err != nil {
			log.Logger.ClientError()
			client.Restart()
		}

		resp := new(detailedvideomodel.POSTResponseDetailedVideoModel)
		resp.Service.Ok = true
		resp.Result.Media = grpcResp.GetMedia()

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Logger.Fatalln(err)
		}
	})

}

func NewDetailedVideoHandler() IDetailedVideo {
	return new(detailedvideo)
}
