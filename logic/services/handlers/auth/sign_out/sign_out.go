package handlers

import (
	signoutmodel "NewPhotoWeb/logic/services/models/auth/sign_out"
	"encoding/json"
	"net/http"
	. "NewPhotoWeb/config"
)

type ISignOutPage interface {
	GetHandler() http.Handler
}

type signout struct{}

func (a *signout) GetHandler() http.Handler {
	//Get handler for account page ...

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := new(signoutmodel.GETResponseSignOutModel)

		http.SetCookie(w, &http.Cookie{Name: "sessionid", Value: "stub", MaxAge: -1, Path: "/"})
		
		resp.Service.Ok = true

		if err := json.NewEncoder(w).Encode(resp); err != nil{
			Logger.Fatalln(err)
		}
	})
}

func NewSignOutPageHandler() ISignOutPage {
	return new(signout)
}
