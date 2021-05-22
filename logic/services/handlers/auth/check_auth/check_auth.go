package handlers

import (
	checkauthmodel "NewPhotoWeb/logic/services/models/auth/check_auth"
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
		session, err := Storage.Get(r, "sessionid")
		if err != nil {
			Logger.Warnln(err.Error())
		}

		resp := new(checkauthmodel.GETResponseCheckAuthModel)
		if _, ok := session.Values["loggedin"]; ok {
			resp.Service.Ok = true
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil{
			Logger.Fatalln(err)
		}
	})

}

func NewCheckAuthHandler() ICheckAuth {
	return new(checkauth)
}
