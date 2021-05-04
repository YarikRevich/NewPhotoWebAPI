package middlewares

import (
	"net/http"
	"strings"

	. "NewPhotoWeb/config"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := Storage.Get(r, "sessionid")
		if err != nil {
			Logger.Warnln(err.Error())
		}
		if _, ok := session.Values["userid"]; ok {
			next.ServeHTTP(w, r)
		} else {
			w.Write([]byte(http.ErrBodyNotAllowed.Error()))
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
