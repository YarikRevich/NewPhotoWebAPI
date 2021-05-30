package handlers

import (
	"NewPhotoWeb/logic/proto"
	checkauthmodel "NewPhotoWeb/logic/services/models/auth/check_auth"
	errormodel "NewPhotoWeb/logic/services/models/error"
	"context"
	"encoding/json"
	"net/http"

	"NewPhotoWeb/log"
	"NewPhotoWeb/logic/client"
)

type ICheckAuth interface {
	GetHandler() http.Handler
}

type checkauth struct{}

func (a *checkauth) GetHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		at := r.Header["X-At"]
		lt := r.Header["X-Lt"]

		if len(at) == 0 || len(lt) == 0 {
			errResp := new(errormodel.ERRORAuthModel)
			errResp.Service.Error = errormodel.AUTH_ERROR
			if err := json.NewEncoder(w).Encode(errResp); err != nil {
				log.Logger.Fatalln(err)
			}
			return
		}

		sourceType := r.Header["S-Type"]

		grpcResp, err := client.NewPhotoAuthClient.RetrieveToken(
			context.Background(),
			&proto.RetrieveTokenRequest{
				AccessToken: at[0],
				LoginToken:  lt[0],
				SourceType:  sourceType[0],
			},
		)
		if err != nil {
			log.Logger.ClientError()
			client.Restart()
		}

		resp := new(checkauthmodel.GETResponseCheckAuthModel)
		if grpcResp.GetOk() {
			w.Header().Add("X-At", grpcResp.GetAccessToken())
			w.Header().Add("X-Lt", grpcResp.GetLoginToken())
			resp.Service.Ok = true
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Logger.Fatalln(err)
		}
	})

}

func NewCheckAuthHandler() ICheckAuth {
	return new(checkauth)
}
