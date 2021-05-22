package middlewares

import (
	"encoding/json"
	"net/http"
	"strings"

	. "NewPhotoWeb/config"
	errormodel "NewPhotoWeb/logic/services/models/error"
	"NewPhotoWeb/utils"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := Storage.Get(r, "sessionid")
		if err != nil {
			Logger.Warnln(err.Error())
		}
		if _, ok := session.Values["userid"]; ok || utils.IsAllowed(r.URL.Path) {
			next.ServeHTTP(w, r)
		} else {
			resp := new(errormodel.ERRORAuthModel)
			resp.Service.Error = errormodel.AUTH_ERROR
			if err := json.NewEncoder(w).Encode(resp); err != nil {
				Logger.Fatalln(err)
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
