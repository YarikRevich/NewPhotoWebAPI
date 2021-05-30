package avatar

import (
	"context"
	"encoding/json"
	"net/http"

	"NewPhotoWeb/logic/proto"
	avatarmodel "NewPhotoWeb/logic/services/models/account/avatar"

	"NewPhotoWeb/log"
	"NewPhotoWeb/logic/client"
)

type IAccountAvatarPage interface {
	GetHandler() http.Handler
	PostHandler() http.Handler
}

type avatar struct{}

func (a *avatar) GetHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		at := r.Header["X-At"]
		lt := r.Header["X-Lt"]

		grpcResp, err := client.NewPhotoClient.GetUserAvatar(
			context.Background(),
			&proto.GetUserAvatarRequest{
				AccessToken: at[0],
				LoginToken:  lt[0],
			},
		)
		if err != nil {
			log.Logger.ClientError()
			client.Restart()
		}

		var resp avatarmodel.GETResponseAvatarModel
		if grpcResp.GetOk() {
			resp.Result.Avatar = grpcResp.GetAvatar()
			resp.Service.Ok = true
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Logger.Fatalln(err)
		}
	})
}

func (a *avatar) PostHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		at := r.Header["X-At"]
		lt := r.Header["X-Lt"]

		var req avatarmodel.POSTRequestAvatarModel
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Logger.Fatalln(err)
		}

		grpcResp, err := client.NewPhotoClient.SetUserAvatar(
			context.Background(),
			&proto.SetUserAvatarRequest{
				AccessToken: at[0],
				LoginToken:  lt[0],
				Avatar:      req.Data.Avatar,
			},
		)
		if err != nil {
			log.Logger.ClientError()
			client.Restart()
		}

		resp := new(avatarmodel.POSTResponseAvatarModel)
		if grpcResp.GetOk() {
			resp.Service.Ok = true
		}

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Logger.Fatalln(err)
		}
	})
}

func NewAvatarHandler() IAccountAvatarPage {
	return new(avatar)
}
