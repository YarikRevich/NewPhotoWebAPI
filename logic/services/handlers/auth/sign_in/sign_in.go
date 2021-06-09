package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"NewPhotoWeb/logic/proto"
	signinmodel "NewPhotoWeb/logic/services/models/auth/sign_in"

	"NewPhotoWeb/log"
	"NewPhotoWeb/logic/client"
)

type ISignInPage interface {
	PostHandler() http.Handler
}

type signin struct{}

func (a *signin) PostHandler() http.Handler {
	//Post handler for account page ...

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req signinmodel.GETRequestSignInModel
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Logger.Fatalln(err)
		}

		sourceType := r.Header["S-Type"]
		var sT proto.SourceType
		switch sourceType[0] {
		case "0":
			sT = proto.SourceType_Web
		case "1":
			sT = proto.SourceType_Mobile
		}

		grpcResp, err := client.NewPhotoAuthClient.LoginUser(
			context.Background(),
			&proto.UserLoginRequest{
				Login:      req.Data.Login,
				Password:   req.Data.Password,
				SourceType: sT,
			})
		if err != nil {
			log.Logger.ClientError()
			client.Restart()
		}

		resp := new(signinmodel.GETResponseSignInModel)

		if grpcResp.GetOk() {
			w.Header().Add("X-At", grpcResp.GetAccessToken())
			w.Header().Add("X-Lt", grpcResp.GetLoginToken())
			resp.Service.Ok = true
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Logger.Fatalln(err)
		}
	})

}

func NewSignInPageHandler() ISignInPage {
	return new(signin)
}
