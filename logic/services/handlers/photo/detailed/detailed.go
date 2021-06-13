package detailed

import (
	"context"
	"encoding/json"
	"net/http"

	"NewPhotoWeb/logic/proto"
	detailedphotomodel "NewPhotoWeb/logic/services/models/photo/detailed"

	"NewPhotoWeb/log"
	"NewPhotoWeb/logic/client"
)

type IDetailedPhoto interface {
	PostHandler() http.Handler
}

type detailedphoto struct{}

func (a *detailedphoto) PostHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var stat detailedphotomodel.POSTRequestDetailedPhotoModel
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
				MediaType:   proto.MediaType_Photo,
			},
		)
		if err != nil {
			log.Logger.ClientError()
			client.Restart()
		}

		resp := new(detailedphotomodel.POSTResponseDetailedPhotoModel)
		resp.Service.Ok = true
		resp.Result.Media = grpcResp.GetMedia()

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Logger.Fatalln(err)
		}
	})

}

func NewDetailedPhotoHandler() IDetailedPhoto {
	return new(detailedphoto)
}
