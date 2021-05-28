package main

import (
	"net/http"
	"os"

	"NewPhotoWeb/log"
	"NewPhotoWeb/logic/services/router"
)

func main() {
	addr, ok := os.LookupEnv("runAddr")
	if !ok {
		log.Logger.Fatalln("runAddr env is not written")
	}

	log.Logger.Fatalln(http.ListenAndServe(addr, router.GetHandler()))
}
