package avatar

import (
	"context"
	"encoding/json"
	"net/http"

	"NewPhotoWeb/logic/proto"
	avatarmodel "NewPhotoWeb/logic/services/models/account/avatar"

	. "NewPhotoWeb/config"
)

type IAccountAvatarPage interface {
	GetHandler() http.Handler
	PostHandler() http.Handler
}

type avatar struct{}

func (a *avatar) GetHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		at, _ := r.Cookie("at")
		lt, _ := r.Cookie("lt")
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
		at, _ := r.Cookie("at")
		lt, _ := r.Cookie("lt")

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
