package handlers

import (
	"math"
	"context"
	"net/http"
	"encoding/json"

	"NewPhotoWeb/logic/proto"
	accountmodel "NewPhotoWeb/logic/services/models/account"
	
	. "NewPhotoWeb/config"
)

type IAccountPage interface {
	GetHandler() http.Handler
}

type account struct{}

func (a *account) GetHandler() http.Handler {
	//Get handler for account page  ...

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := Storage.Get(r, "sessionid")
		if err != nil {
			Logger.Warn(err.Error())
		}

		grpcResp, err := NPC.GetUserinfo(context.Background(), &proto.GetUserinfoRequest{Userid: session.Values["userid"].(string)})
		if err != nil {
			Logger.ClientError()
		}

		resp := new(accountmodel.GETResponseAccountModel)

		resp.Result.Firstname = grpcResp.GetFirstname()
		resp.Result.Secondname = grpcResp.GetSecondname()
		resp.Result.Storage = math.Round((grpcResp.GetStorage()*math.Pow(10, -9))*100) / 100
		resp.Service.Ok = true

		if err := json.NewEncoder(w).Encode(resp); err != nil{
			Logger.Fatalln(err)
		}
	})
}

func NewAccountPageHandler() IAccountPage {
	return new(account)
}
