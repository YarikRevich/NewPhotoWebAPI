package avatar

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"NewPhotoWeb/logic/proto"
	avatarmodel "NewPhotoWeb/logic/services/models/account/avatar"
	errormodel "NewPhotoWeb/logic/services/models/error"

	. "NewPhotoWeb/config"
)

type IAccountAvatarPage interface {
	GetHandler() http.Handler
	PostHandler() http.Handler
}

type avatar struct{}

func (a *avatar) GetHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Cookies())
		fmt.Println(r.Cookie("at"))
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
		grpcResp, err := NPC.GetUserAvatar(
			context.Background(),
			&proto.GetUserAvatarRequest{
				AccessToken: at.Value,
				LoginToken:  lt.Value,
			},
		)
		if err != nil {
			Logger.ClientError()
		}

		var resp avatarmodel.GETResponseAvatarModel
		if grpcResp.GetOk() {
			resp.Result.Avatar = string(grpcResp.GetAvatar())
			resp.Service.Ok = true
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			Logger.Fatalln(err)
		}
	})
}

func (a *avatar) PostHandler() http.Handler {
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

		var req avatarmodel.POSTRequestAvatarModel
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			Logger.Fatalln(err)
		}

		grpcResp, err := NPC.SetUserAvatar(
			context.Background(),
			&proto.SetUserAvatarRequest{
				AccessToken: at.Value,
				LoginToken:  lt.Value,
				Avatar:      []byte(req.Data.Avatar),
			},
		)
		if err != nil {
			Logger.ClientError()
		}

		resp := new(avatarmodel.POSTResponseAvatarModel)
		if grpcResp.GetOk() {
			resp.Service.Ok = true
		}

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			Logger.Fatalln(err)
		}
	})
}

func NewAvatarHandler() IAccountAvatarPage {
	return new(avatar)
}
