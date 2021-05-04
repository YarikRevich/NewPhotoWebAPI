package handlers

import (
	// "NewPhotoWeb/logic/proto"
	// "NewPhotoWeb/logic/services/models"
	// "context"
	// "encoding/json"

	"fmt"
	"net/http"
	"os"

	. "NewPhotoWeb/config"
)

type IRetrieveTokenPage interface {
	GetHandler() http.Handler
}

type retrievetoken struct{}

func (b *retrievetoken) GetHandler() http.Handler {
	//Get handler for retrieve_token ...

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := Storage.Get(r, "sessionid")
		if err != nil {
			err = session.Save(r, w)
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				return
			}
			return
		}
	})
}

func NewRetrieveTokenHandler() IRetrieveTokenPage {
	return new(retrievetoken)
}
