package handlers

import (
	"context"
	"encoding/json"
	"math"
	"net/http"

	"NewPhotoWeb/logic/proto"
	accountmodel "NewPhotoWeb/logic/services/models/account"
	errormodel "NewPhotoWeb/logic/services/models/error"

	. "NewPhotoWeb/config"
)

type IAccountPage interface {
	GetHandler() http.Handler
}

type account struct{}

func (a *account) GetHandler() http.Handler {
	//Get handler for account page  ...

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errResp := new(errormodel.ERRORAuthModel)
		errResp.Service.Error = errormodel.AUTH_ERROR
		at, err := r.Cookie("at")
		if err != nil {
			if err := json.NewEncoder(w).Encode(errResp); err != nil {
				Logger.Fatalln(err)
			}
		}
		lt, err := r.Cookie("lt")
		if err != nil {
			if err := json.NewEncoder(w).Encode(errResp); err != nil {
				Logger.Fatalln(err)
			}
		}

		grpcResp, err := NPC.GetUserinfo(
			context.Background(),
			&proto.GetUserinfoRequest{
				AccessToken: at.Value,
				LoginToken:  lt.Value,
			},
		)
		if err != nil {
			Logger.ClientError()
		}

		resp := new(accountmodel.GETResponseAccountModel)

		resp.Result.Firstname = grpcResp.GetFirstname()
		resp.Result.Secondname = grpcResp.GetSecondname()
		resp.Result.Storage = math.Round((grpcResp.GetStorage()*math.Pow(10, -9))*100) / 100
		resp.Service.Ok = true

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			Logger.Fatalln(err)
		}
	})
}

func NewAccountPageHandler() IAccountPage {
	return new(account)
}
