package middlewares

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	. "NewPhotoWeb/config"
	"NewPhotoWeb/logic/proto"
	errormodel "NewPhotoWeb/logic/services/models/error"
	"NewPhotoWeb/utils"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errResp := new(errormodel.ERRORAuthModel)
		errResp.Service.Error = errormodel.AUTH_ERROR
		at, err := r.Cookie("at")
		if err != nil {
			if err := json.NewEncoder(w).Encode(errResp); err != nil {
				Logger.Fatalln(err)
			}
		}
		lt, err := r.Cookie("lt")
		if err != nil {
			if err := json.NewEncoder(w).Encode(errResp); err != nil {
				Logger.Fatalln(err)
			}
		}
		if utils.IsAllowed(r.URL.Path) {
			next.ServeHTTP(w, r)
		} else {
			grpcResp, err := AC.RetrieveToken(context.Background(), &proto.RetrieveTokenRequest{AccessToken: at.Value, LoginToken: lt.Value})
			if err != nil {
				Logger.ClientError()
			}
			if grpcResp.GetOk() {
				http.SetCookie(w, &http.Cookie{Name: "at", Value: grpcResp.AccessToken})
				http.SetCookie(w, &http.Cookie{Name: "lt", Value: grpcResp.LoginToken})
				next.ServeHTTP(w, r)
			} else {
				resp := new(errormodel.ERRORAuthModel)
				resp.Service.Error = errormodel.NOT_THIS_TIME_ERROR
				if err := json.NewEncoder(w).Encode(resp); err != nil {
					Logger.Fatalln(err)
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
