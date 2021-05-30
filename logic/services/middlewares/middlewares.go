package middlewares

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"NewPhotoWeb/log"
	"NewPhotoWeb/logic/client"
	"NewPhotoWeb/logic/proto"
	errormodel "NewPhotoWeb/logic/services/models/error"
	"NewPhotoWeb/utils"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if utils.IsAllowed(r.URL.Path) {
			next.ServeHTTP(w, r)
		} else {
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

			grpcResp, err := client.NewPhotoAuthClient.IsTokenCorrect(
				context.Background(),
				&proto.IsTokenCorrectRequest{
					AccessToken: at[0],
					LoginToken:  lt[0],
					SourceType:  sourceType[0],
				},
			)
			if err != nil {
				log.Logger.ClientError()
				client.Restart()
			}

			if grpcResp.GetOk() {
				next.ServeHTTP(w, r)
			} else {
				resp := new(errormodel.ERRORAuthModel)
				resp.Service.Error = errormodel.NOT_THIS_TIME_ERROR
				if err := json.NewEncoder(w).Encode(resp); err != nil {
					log.Logger.Fatalln(err)
				}
			}
		}
	})
}

func FetchingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fetch := r.Header.Get("Fetch")
		if strings.Contains(r.URL.Path, "spa") && fetch == "true" || !strings.Contains(r.URL.Path, "spa") {
			next.ServeHTTP(w, r)
		} else {
			http.NotFound(w, r)
		}
	})
}

func EnableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", r.URL.Host)
		next.ServeHTTP(w, r)
	})
}
