package detailed

import (
	"context"
	"net/http"
	"encoding/json"
	"encoding/base64"

	"NewPhotoWeb/logic/proto"
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

		session, err := Storage.Get(r, "sessionid")
		if err != nil {
			Logger.Warnln(err.Error())
		}

		grpcResp, err := NPC.GetFullPhotoByThumbnail(context.Background(), &proto.GetFullPhotoByThumbnailRequest{Userid: session.Values["userid"].(string), Thumbnail: efile})
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
