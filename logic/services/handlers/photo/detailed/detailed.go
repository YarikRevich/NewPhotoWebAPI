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
		var req detailedphotomodel.GETRequestDetailedPhotoModel
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Logger.Fatalln(err)
		}

		at := r.Header["X-At"]
		lt := r.Header["X-Lt"]

		grpcResp, err := client.NewPhotoClient.GetFullPhotoByThumbnail(
			context.Background(),
			&proto.GetFullPhotoByThumbnailRequest{
				AccessToken: at[0],
				LoginToken:  lt[0],
				Thumbnail:   req.Data.Thumbnail,
			},
		)
		if err != nil {
			log.Logger.ClientError(); client.Restart()
		}

		resp := new(detailedphotomodel.GETResponseDetailedPhotoModel)
		resp.Service.Ok = true
		resp.Result.Photo = grpcResp.GetPhoto()

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Logger.Fatalln(err)
		}
	})

}

func NewDetailedPhotoHandler() IDetailedPhoto {
	return new(detailedphoto)
}
