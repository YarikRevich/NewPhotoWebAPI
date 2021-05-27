package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"NewPhotoWeb/logic/proto"
	signinmodel "NewPhotoWeb/logic/services/models/auth/sign_in"

	. "NewPhotoWeb/config"
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
			Logger.Fatalln(err)
		}

		grpcResp, err := AC.LoginUser(
			context.Background(),
			&proto.UserLoginRequest{
				Login:    req.Data.Login,
				Password: req.Data.Password,
			})
		if err != nil {
			Logger.ClientError()
		}

		resp := new(signinmodel.GETResponseSignInModel)

		if grpcResp.GetOk() {
			http.SetCookie(w, &http.Cookie{Name: "at", Value: grpcResp.AccessToken})
			http.SetCookie(w, &http.Cookie{Name: "lt", Value: grpcResp.LoginToken})
			resp.Service.Ok = true
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			Logger.Fatalln(err)
		}
	})

}

func NewSignInPageHandler() ISignInPage {
	return new(signin)
}
