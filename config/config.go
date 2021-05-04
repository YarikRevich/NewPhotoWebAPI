package config

import(
	"NewPhotoWeb/log"
	"NewPhotoWeb/logic/client"

	"github.com/gorilla/sessions"
	"github.com/gorilla/securecookie"
)

var (
	//Session storage
	Storage = sessions.NewCookieStore(securecookie.GenerateRandomKey(64), securecookie.GenerateRandomKey(32)) 

	//New Photo client
	NPC  = client.NewPhotoClient()

	//Authentication client
	AC   = client.NewAuthClient()

	//Logger for web app
	Logger = log.New()
)