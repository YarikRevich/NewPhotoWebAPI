package handlers

import (
	"NewPhotoWeb/logic/proto"
	signupmodel "NewPhotoWeb/logic/services/models/auth/sign_up"
	"context"
	"encoding/json"

	"net/http"


	"NewPhotoWeb/log"
	"NewPhotoWeb/logic/client"
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
			log.Logger.Fatalln(err)
		}

		resp := new(signupmodel.POSTResponseRegestrationModel)

		grpcResp, err := client.NewPhotoAuthClient.RegisterUser(
			context.Background(),
			&proto.UserRegisterRequest{
				Login:      req.Data.Login,
				Password:   req.Data.Password1,
				Firstname:  req.Data.Firstname,
				Secondname: req.Data.Secondname,
			})
		if err != nil {
			log.Logger.ClientError(); client.Restart()
		}

		resp.Service.Ok = grpcResp.GetOk()

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Logger.Fatalln(err)
		}
	})
}

func NewSignUpHandler() ISignUpPage {
	return new(signup)
}
