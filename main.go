package main

import (
	"net/http"
	"os"

	. "NewPhotoWeb/config"
	"NewPhotoWeb/logic/services/router"
)

func main() {
	addr, ok := os.LookupEnv("runAddr")
	if !ok {
		Logger.Fatalln("runAddr env is not written")
	}

	Logger.Fatalln(http.ListenAndServe(addr, router.GetHandler()))
}
