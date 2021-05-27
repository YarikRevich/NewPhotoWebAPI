package handlers

import (
	"NewPhotoWeb/logic/proto"
	checkauthmodel "NewPhotoWeb/logic/services/models/auth/check_auth"
	errormodel "NewPhotoWeb/logic/services/models/error"
	"context"
	"encoding/json"
	"net/http"

	. "NewPhotoWeb/config"
)

type ICheckAuth interface {
	GetHandler() http.Handler
}

type checkauth struct{}

func (a *checkauth) GetHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errResp := new(errormodel.ERRORAuthModel)
		errResp.Service.Error = errormodel.AUTH_ERROR
		at, err := r.Cookie("at")
		if err != nil {
			if err := json.NewEncoder(w).Encode(errResp); err != nil {
				Logger.Fatalln(err)
			}
			return
		}
		lt, err := r.Cookie("lt")
		if err != nil {
			if err := json.NewEncoder(w).Encode(errResp); err != nil {
				Logger.Fatalln(err)
			}
			return
		}

		grpcResp, err := AC.RetrieveToken(context.Background(), &proto.RetrieveTokenRequest{AccessToken: at.Value, LoginToken: lt.Value})
		if err != nil {
			Logger.ClientError()
		}
		resp := new(checkauthmodel.GETResponseCheckAuthModel)
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

func NewCheckAuthHandler() ICheckAuth {
	return new(checkauth)
}
