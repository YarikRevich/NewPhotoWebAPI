package tests

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"

	signupmodel "NewPhotoWeb/logic/services/models/auth/sign_up"
	"NewPhotoWeb/logic/services/router"

	"github.com/franela/goblin"
)

func TestSignUpHandler(t *testing.T) {
	b := goblin.Goblin(t)
	b.Describe("Check auth handler", func() {
		b.It("Should get equal resp", func() {
			s := httptest.NewServer(router.GetHandler())

			randInt, err := rand.Int(rand.Reader, big.NewInt(1000))
			if err != nil {
				b.Fail(err)
			}

			q := signupmodel.POSTRequestRegestrationModel{
				Data: struct {
					Login      string "json:\"login\""
					Firstname  string "json:\"firstname\""
					Secondname string "json:\"secondname\""
					Password1  string "json:\"password1\""
					Password2  string "json:\"password2\""
				}{
					Login:      "Test" + randInt.String(),
					Firstname:  "Test" + randInt.String(),
					Secondname: "Test" + randInt.String(),
					Password1:  "Test" + randInt.String(),
					Password2:  "Test" + randInt.String(),
				},
			}

			by, err := json.Marshal(q)
			if err != nil {
				b.Fail(err)
			}

			rr, err := http.Post(s.URL+router.SignUpPath, "application/json", bytes.NewReader(by))
			if err != nil {
				b.Fail(err)
			}
			a := signupmodel.POSTResponseRegestrationModel{}
			if err := json.NewDecoder(rr.Body).Decode(&a); err != nil {
				b.Fail(err)
			}

			b.Assert(a.Service.Ok).Equal(true, "Registration is failed")
		})
	})
}
