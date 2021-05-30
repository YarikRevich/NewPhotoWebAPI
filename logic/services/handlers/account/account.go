package handlers

import (
	"context"
	"encoding/json"
	"math"
	"net/http"

	"NewPhotoWeb/logic/proto"
	accountmodel "NewPhotoWeb/logic/services/models/account"

	"NewPhotoWeb/log"
	"NewPhotoWeb/logic/client"
)

type IAccountPage interface {
	GetHandler() http.Handler
	DeleteHandler() http.Handler
}

type account struct{}

func (a *account) GetHandler() http.Handler {
	//Get handler for account page  ...

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		at := r.Header["X-At"]
		lt := r.Header["X-Lt"]

		grpcResp, err := client.NewPhotoClient.GetUserinfo(
			context.Background(),
			&proto.GetUserinfoRequest{
				AccessToken: at[0],
				LoginToken:  lt[0],
			},
		)
		if err != nil {
			log.Logger.ClientError(); client.Restart()
		}

		resp := new(accountmodel.GETResponseAccountModel)

		if grpcResp.GetOk() {
			resp.Result.Firstname = grpcResp.GetFirstname()
			resp.Result.Secondname = grpcResp.GetSecondname()
			resp.Result.Storage = math.Round((grpcResp.GetStorage()*math.Pow(10, -9))*100) / 100
			resp.Service.Ok = true
		}

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Logger.Fatalln(err)
		}
	})
}

func (a *account) DeleteHandler() http.Handler {
	//Get handler for account page  ...

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		at := r.Header["X-At"]
		lt := r.Header["X-Lt"]

		grpcResp, err := client.NewPhotoClient.DeleteAccount(
			context.Background(),
			&proto.DeleteAccountRequest{
				AccessToken: at[0],
				LoginToken:  lt[0],
			},
		)
		if err != nil {
			log.Logger.ClientError(); client.Restart()
		}

		resp := new(accountmodel.DELETEResponseAccountModel)
		if grpcResp.GetOk() {
			resp.Service.Ok = true
		}

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Logger.Fatalln(err)
		}
	})
}

func NewAccountPageHandler() IAccountPage {
	return new(account)
}
