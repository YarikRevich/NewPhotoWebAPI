package main

import (
	"net/http"
	"os"

	account "NewPhotoWeb/logic/services/handlers/account"
	album "NewPhotoWeb/logic/services/handlers/album"
	checkauth "NewPhotoWeb/logic/services/handlers/auth/check_auth"
	signin "NewPhotoWeb/logic/services/handlers/auth/sign_in"
	signout "NewPhotoWeb/logic/services/handlers/auth/sign_out"
	signup "NewPhotoWeb/logic/services/handlers/auth/sign_up"
	photo "NewPhotoWeb/logic/services/handlers/photo"
	detailedphoto "NewPhotoWeb/logic/services/handlers/photo/detailed"
	"NewPhotoWeb/logic/services/middlewares"

	. "NewPhotoWeb/config"

	"github.com/gorilla/mux"
)

var (
	router = mux.NewRouter()
)

func init() {
	router.Use(middlewares.FetchingMiddleware)
	router.Use(middlewares.EnableCORS)

	router.HandleFunc("/check_auth", func(w http.ResponseWriter, r *http.Request) {
		pagehandler := checkauth.NewCheckAuthHandler()
		pagehandler.GetHandler().ServeHTTP(w, r)
	}).Methods("GET")

	router.HandleFunc("/sign_up", func(w http.ResponseWriter, r *http.Request) {
		pagehandler := signup.NewSignUpHandler()
		pagehandler.PostHandler().ServeHTTP(w, r)
	}).Methods("POST")

	router.HandleFunc("/sign_in", func(w http.ResponseWriter, r *http.Request) {
		pagehandler := signin.NewSignInPageHandler()
		pagehandler.PostHandler().ServeHTTP(w, r)
	}).Methods("POST")

	router.HandleFunc("/sign_out", func(w http.ResponseWriter, r *http.Request) {
		pagehandler := signout.NewSignOutPageHandler()
		pagehandler.GetHandler().ServeHTTP(w, r)
	}).Methods("GET")

	router.HandleFunc("/photos/detailed", func(w http.ResponseWriter, r *http.Request) {
		pagehandler := detailedphoto.NewDetailedPhotoHandler()
		pagehandler.GetHandler().ServeHTTP(w, r)
	}).Methods("GET")

	router.HandleFunc("/photos", func(w http.ResponseWriter, r *http.Request) {
		pagehandler := photo.NewPhotoPageHandler()
		switch r.Method {
		case "GET":
			pagehandler.GetHandler().ServeHTTP(w, r)
		case "POST":
			pagehandler.PostHandler().ServeHTTP(w, r)
		}
	}).Methods("GET", "POST")

	router.HandleFunc("/albums", func(w http.ResponseWriter, r *http.Request){
		pagehandler := album.NewAlbumHandler()
		switch r.Method{
		case "GET":
			pagehandler.GetHandler().ServeHTTP(w, r)
		case "POST":
			pagehandler.PostHandler().ServeHTTP(w, r)
		case "PUT":
			pagehandler.PutHandler().ServeHTTP(w, r)
		case "DELETE":
			pagehandler.DeleteHandler().ServeHTTP(w, r)
		}
	}).Methods("GET", "POST", "PUT", "DELETE")

	router.HandleFunc("/account", func(w http.ResponseWriter, r *http.Request) {
		pagehandler := account.NewAccountPageHandler()
		pagehandler.GetHandler().ServeHTTP(w, r)
	}).Methods("GET")
}

func main() {
	addr, ok := os.LookupEnv("runAddr")
	if !ok{
		Logger.Fatalln("runAddr env is not written")
	}

	Logger.Fatalln(http.ListenAndServe(addr, router))
}
