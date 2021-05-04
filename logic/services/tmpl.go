package services

import (
	// "NewPhotoWeb/utils"
	"html/template"
	"net/http"
	// . "NewPhotoWeb/config"
)

func Index(w http.ResponseWriter, r *http.Request){
	t := template.Must(template.ParseFiles("templates/app/build/index.html"))
	t.Execute(w, nil)	
}

// func AlbumsPage(w http.ResponseWriter, r *http.Request, t utils.AlbumPageVars) {
// 	temp, err := template.ParseFiles("templates/albums.tmpl", "templates/index.tmpl")
// 	if err != nil {
// 		Logger.Fatalln(err.Error())
// 	}
// 	err = temp.Execute(w, t)
// 	if err != nil{
// 		Logger.Fatalln(err.Error())
// 	}
// }

// func AccountPage(w http.ResponseWriter, r *http.Request, tempres *utils.AccountPageVars) {
// 	temp, err := template.ParseFiles("templates/account.tmpl", "templates/index.tmpl")
// 	if err != nil {
// 		Logger.Fatalln(err.Error())
// 	}
// 	err = temp.Execute(w, tempres)
// 	if err != nil {
// 		Logger.Fatalln(err.Error())
// 	}
// }

// func PhotosPage(w http.ResponseWriter, r *http.Request, tempres *utils.PhotoPageVars) {
// 	temp, err := template.ParseFiles("templates/photos.tmpl", "templates/index.tmpl")
// 	if err != nil {
// 		Logger.Fatalln(err.Error())
// 	}
// 	err = temp.Execute(w, tempres)
// 	if err != nil {
// 		Logger.Fatalln(err.Error())
// 	}
// }

// func WelcomePage(w http.ResponseWriter, r *http.Request) {
// 	temp, err := template.ParseFiles("templates/welcome.tmpl", "templates/index.tmpl")
// 	if err != nil {
// 		Logger.Fatalln(err.Error())
// 	}
// 	err = temp.Execute(w, nil)
// 	if err != nil {
// 		Logger.Fatalln(err.Error())
// 	}

// }

// func RegistrationPage(w http.ResponseWriter, r *http.Request, t *utils.RegPageVars) {
// 	temp, err := template.ParseFiles("templates/registration.tmpl", "templates/index.tmpl")
// 	if err != nil {
// 		Logger.Fatalln(err.Error())
// 	}
// 	err = temp.Execute(w, t)
// 	if err != nil{
// 		Logger.Fatalln(err.Error())
// 	}
// }

// func EqualAlbumPage(w http.ResponseWriter, r *http.Request, t *utils.EqualAlbumVars) {
// 	temp, err := template.ParseFiles("templates/equalalbum.tmpl", "templates/index.tmpl")
// 	if err != nil {
// 		Logger.Fatalln(err.Error())
// 	}
// 	err = temp.Execute(w, t)
// 	if err != nil{
// 		Logger.Fatalln(err.Error())
// 	}
// }
