package handlers

import (
	"context"
	"encoding/json"
	"math"
	"net/http"

	"NewPhotoWeb/logic/proto"
	accountmodel "NewPhotoWeb/logic/services/models/account"

	. "NewPhotoWeb/config"
)

type IAccountPage interface {
	GetHandler() http.Handler
	DeleteHandler() http.Handler
}

type account struct{}

func (a *account) GetHandler() http.Handler {
	//Get handler for account page  ...

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		at, _ := r.Cookie("at")
		lt, _ := r.Cookie("lt")

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

		if grpcResp.GetOk() {
			resp.Result.Firstname = grpcResp.GetFirstname()
			resp.Result.Secondname = grpcResp.GetSecondname()
			resp.Result.Storage = math.Round((grpcResp.GetStorage()*math.Pow(10, -9))*100) / 100
			resp.Service.Ok = true
		}

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			Logger.Fatalln(err)
		}
	})
}

func (a *account) DeleteHandler() http.Handler {
	//Get handler for account page  ...

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		at, _ := r.Cookie("at")
		lt, _ := r.Cookie("lt")

		grpcResp, err := NPC.DeleteAccount(
			context.Background(),
			&proto.DeleteAccountRequest{
				AccessToken: at.Value,
				LoginToken:  lt.Value,
			},
		)
		if err != nil {
			Logger.ClientError()
		}

		resp := new(accountmodel.DELETEResponseAccountModel)
		if grpcResp.GetOk() {
			resp.Service.Ok = true
		}

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			Logger.Fatalln(err)
		}
	})
}

func NewAccountPageHandler() IAccountPage {
	return new(account)
}
