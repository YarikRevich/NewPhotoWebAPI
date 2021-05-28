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

		grpcResp, err := client.NewPhotoAuthClient.LoginUser(
			context.Background(),
			&proto.UserLoginRequest{
				Login:    req.Data.Login,
				Password: req.Data.Password,
			})
		if err != nil {
			log.Logger.ClientError()
			client.Restart()
		}

		resp := new(signinmodel.GETResponseSignInModel)

		if grpcResp.GetOk() {
			http.SetCookie(w, &http.Cookie{Name: "at", Value: grpcResp.AccessToken, Path: "/"})
			http.SetCookie(w, &http.Cookie{Name: "lt", Value: grpcResp.LoginToken, Path: "/"})
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
