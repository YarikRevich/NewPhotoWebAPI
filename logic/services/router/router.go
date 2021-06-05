package router

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"

	account "NewPhotoWeb/logic/services/handlers/account"
	avatar "NewPhotoWeb/logic/services/handlers/account/avatar"
	album "NewPhotoWeb/logic/services/handlers/album"
	detailedalbum "NewPhotoWeb/logic/services/handlers/album/detailed"
	infodetailedalbum "NewPhotoWeb/logic/services/handlers/album/detailed/info"
	checkauth "NewPhotoWeb/logic/services/handlers/auth/check_auth"
	signin "NewPhotoWeb/logic/services/handlers/auth/sign_in"
	signout "NewPhotoWeb/logic/services/handlers/auth/sign_out"
	signup "NewPhotoWeb/logic/services/handlers/auth/sign_up"
	photo "NewPhotoWeb/logic/services/handlers/photo"
	detailedphoto "NewPhotoWeb/logic/services/handlers/photo/detailed"
	detailedvideo "NewPhotoWeb/logic/services/handlers/video/detailed"
	video "NewPhotoWeb/logic/services/handlers/video"
	"NewPhotoWeb/logic/services/middlewares"
)

const (
	CheckAuthPath          = "/check_auth"
	SignUpPath             = "/sign_up"
	SignInPath             = "/sign_in"
	SignOutPath            = "/sign_out"
	PhotosDetailedPath     = "/photos/detailed"
	PhotosPath             = "/photos"
	VideosDetailedPath     = "/videos/detailed"
	VideosPath             = "/videos"
	AlbumsPath             = "/albums"
	AlbumsDetailedPath     = "/albums/detailed"
	AlbumsDetailedInfoPath = "/albums/detailed/info"
	AccountPath            = "/account"
	AccountAvatarPath      = "/account/avatar"
)

func GetHandler() *mux.Router {
	router := mux.NewRouter()

	router.Use(middlewares.FetchingMiddleware)
	router.Use(middlewares.EnableCORS)
	router.Use(middlewares.AuthMiddleware)

	router.HandleFunc(CheckAuthPath, func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(w.Header())
		pagehandler := checkauth.NewCheckAuthHandler()
		pagehandler.GetHandler().ServeHTTP(w, r)
	}).Methods("GET", "OPTIONS")

	router.HandleFunc(SignUpPath, func(w http.ResponseWriter, r *http.Request) {
		pagehandler := signup.NewSignUpHandler()
		pagehandler.PostHandler().ServeHTTP(w, r)
	}).Methods("POST", "OPTIONS")

	router.HandleFunc(SignInPath, func(w http.ResponseWriter, r *http.Request) {
		pagehandler := signin.NewSignInPageHandler()
		pagehandler.PostHandler().ServeHTTP(w, r)
	}).Methods("POST", "OPTIONS")

	router.HandleFunc(SignOutPath, func(w http.ResponseWriter, r *http.Request) {
		pagehandler := signout.NewSignOutPageHandler()
		pagehandler.GetHandler().ServeHTTP(w, r)
	}).Methods("GET", "OPTIONS")

	router.HandleFunc(PhotosDetailedPath, func(w http.ResponseWriter, r *http.Request) {
		pagehandler := detailedphoto.NewDetailedPhotoHandler()
		pagehandler.GetHandler().ServeHTTP(w, r)
	}).Methods("GET", "OPTIONS")

	router.HandleFunc(PhotosPath, func(w http.ResponseWriter, r *http.Request) {
		pagehandler := photo.NewPhotoPageHandler()
		switch r.Method {
		case "GET":
			pagehandler.GetHandler().ServeHTTP(w, r)
		case "POST":
			pagehandler.PostHandler().ServeHTTP(w, r)
		}
	}).Methods("GET", "POST", "OPTIONS")

	router.HandleFunc(VideosPath, func(w http.ResponseWriter, r *http.Request) {
		pagehandler := video.NewVideoPageHandler()
		pagehandler.PostHandler().ServeHTTP(w, r)
	}).Methods("POST", "OPTIONS")

	router.HandleFunc(VideosDetailedPath, func(w http.ResponseWriter, r *http.Request) {
		pagehandler := detailedvideo.NewDetailedVideoHandler()
		pagehandler.GetHandler().ServeHTTP(w, r)
	}).Methods("GET",  "OPTIONS")

	router.HandleFunc(AlbumsPath, func(w http.ResponseWriter, r *http.Request) {
		pagehandler := album.NewAlbumHandler()
		switch r.Method {
		case "GET":
			pagehandler.GetHandler().ServeHTTP(w, r)
		case "POST":
			pagehandler.PostHandler().ServeHTTP(w, r)
		case "DELETE":
			pagehandler.DeleteHandler().ServeHTTP(w, r)
		}
	}).Methods("GET", "POST", "DELETE", "OPTIONS")

	router.HandleFunc(AlbumsDetailedPath, func(w http.ResponseWriter, r *http.Request) {
		pagehandler := detailedalbum.NewDetailedAlbumPageHandler()
		switch r.Method {
		case "GET":
			pagehandler.GetHandler().ServeHTTP(w, r)
		case "PUT":
			pagehandler.PutHandler().ServeHTTP(w, r)
		case "DELETE":
			pagehandler.DeleteHandler().ServeHTTP(w, r)
		}
	}).Methods("GET", "PUT", "DELETE", "OPTIONS")

	router.HandleFunc(AlbumsDetailedInfoPath, func(w http.ResponseWriter, r *http.Request) {
		pagehandler := infodetailedalbum.NewInfoDetailedAlbumPageHandler()
		pagehandler.GetHandler().ServeHTTP(w, r)
	})

	router.HandleFunc(AccountPath, func(w http.ResponseWriter, r *http.Request) {
		pagehandler := account.NewAccountPageHandler()
		switch r.Method {
		case "GET":
			pagehandler.GetHandler().ServeHTTP(w, r)
		case "DELETE":
			pagehandler.DeleteHandler().ServeHTTP(w, r)
		}
	}).Methods("GET", "DELETE", "OPTIONS")

	router.HandleFunc(AccountAvatarPath, func(w http.ResponseWriter, r *http.Request) {
		pagehandler := avatar.NewAvatarHandler()
		switch r.Method {
		case "GET":
			pagehandler.GetHandler().ServeHTTP(w, r)
		case "POST":
			pagehandler.PostHandler().ServeHTTP(w, r)
		}
	}).Methods("GET", "POST", "OPTIONS")

	return router
}
