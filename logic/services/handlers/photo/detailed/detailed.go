package detailed

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"

	"NewPhotoWeb/logic/proto"
	errormodel "NewPhotoWeb/logic/services/models/error"
	detailedphotomodel "NewPhotoWeb/logic/services/models/photo/detailed"

	. "NewPhotoWeb/config"
)

type IDetailedPhoto interface {
	GetHandler() http.Handler
}

type detailedphoto struct{}

func (a *detailedphoto) GetHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req detailedphotomodel.GETRequestDetailedPhotoModel
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			Logger.Fatalln(err)
		}

		efile, err := base64.StdEncoding.DecodeString(req.Data.Photo)
		if err != nil {
			Logger.Fatalln(err)
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

		grpcResp, err := NPC.GetFullPhotoByThumbnail(
			context.Background(),
			&proto.GetFullPhotoByThumbnailRequest{
				AccessToken: at.Value,
				LoginToken:  lt.Value,
				Thumbnail:   efile,
			},
		)
		if err != nil {
			Logger.Fatalln(err)
		}

		resp := new(detailedphotomodel.GETResponseDetailedPhotoModel)
		resp.Service.Ok = true
		resp.Result.Photo = base64.StdEncoding.EncodeToString(grpcResp.GetPhoto())

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			Logger.Fatalln(err)
		}
	})

}

func NewDetailedPhotoHandler() IDetailedPhoto {
	return new(detailedphoto)
}
