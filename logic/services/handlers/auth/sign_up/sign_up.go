package handlers

import (
	"NewPhotoWeb/logic/proto"
	signupmodel "NewPhotoWeb/logic/services/models/auth/sign_up"
	"context"
	"encoding/json"

	"net/http"

	. "NewPhotoWeb/config"
)

type ISignUpPage interface {
	PostHandler() http.Handler
}

type signup struct{}

func (a *signup) PostHandler() http.Handler {
	//Post handler for account page ...

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req signupmodel.POSTRequestRegestrationModel
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			Logger.Fatalln(err)
		}

		resp := new(signupmodel.POSTResponseRegestrationModel)

		grpcResp, err := AC.RegisterUser(
			context.Background(),
			&proto.UserRegisterRequest{
				Login:      req.Data.Login,
				Password:   req.Data.Password1,
				Firstname:  req.Data.Firstname,
				Secondname: req.Data.Secondname,
			})
		if err != nil {
			Logger.ClientError()
		}

		resp.Service.Ok = grpcResp.GetOk()

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			Logger.Fatalln(err)
		}
	})
}

func NewSignUpHandler() ISignUpPage {
	return new(signup)
}
