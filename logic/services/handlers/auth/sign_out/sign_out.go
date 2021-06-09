package handlers

import (
	"NewPhotoWeb/logic/client"
	"NewPhotoWeb/logic/proto"
	signoutmodel "NewPhotoWeb/logic/services/models/auth/sign_out"
	"context"
	"encoding/json"
	"net/http"

	"NewPhotoWeb/log"
)

type ISignOutPage interface {
	DeleteHandler() http.Handler
}

type signout struct{}

func (a *signout) DeleteHandler() http.Handler {
	//Get handler for account page ...

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := new(signoutmodel.GETResponseSignOutModel)

		sourceType := r.Header["S-Type"]
		var sT proto.SourceType
		switch sourceType[0] {
		case "0":
			sT = proto.SourceType_Web
		case "1":
			sT = proto.SourceType_Mobile
		}

		grpcResp, err := client.NewPhotoAuthClient.LogoutUser(
			context.Background(),
			&proto.UserLogoutRequest{
				AccessToken: r.Header.Get("X-At"),
				LoginToken: r.Header.Get("X-Lt"),
				SourceType:  sT,
			})
		if err != nil {
			log.Logger.ClientError()
			client.Restart()
		}
		if grpcResp.GetOk() {
			resp.Service.Ok = true
		}

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Logger.Fatalln(err)
		}
	})
}

func NewSignOutPageHandler() ISignOutPage {
	return new(signout)
}
