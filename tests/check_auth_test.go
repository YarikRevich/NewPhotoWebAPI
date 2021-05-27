package tests

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"

	checkauthmodel "NewPhotoWeb/logic/services/models/auth/check_auth"
	signinmodel "NewPhotoWeb/logic/services/models/auth/sign_in"
	signupmodel "NewPhotoWeb/logic/services/models/auth/sign_up"
	"NewPhotoWeb/logic/services/router"

	"github.com/franela/goblin"
)

//Tests check_auth handler, firstly it regs new user
//Then it logins and after that checks if user is authed
func TestCheckAuthHandler(t *testing.T) {
	b := goblin.Goblin(t)
	b.Describe("Check auth handler", func() {
		b.It("Should get equal resp", func() {
			s := httptest.NewServer(router.GetHandler())

			jar, err := cookiejar.New(nil)
			if err != nil {
				b.Fail(err)
			}
			c := http.DefaultClient
			c.Jar = jar
					

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

			rr, err := c.Post(s.URL+router.SignUpPath, "application/json", bytes.NewReader(by))
			if err != nil {
				b.Fail(err)
			}
			a := signupmodel.POSTResponseRegestrationModel{}
			if err := json.NewDecoder(rr.Body).Decode(&a); err != nil {
				b.Fail(err)
			}

			b.Assert(a.Service.Ok).Equal(true, "Registration is failed")

			t := signinmodel.GETRequestSignInModel{
				Data: struct {
					Login    string "json:\"login\""
					Password string "json:\"password\""
				}{
					Login:    "Test" + randInt.String(),
					Password: "Test" + randInt.String(),
				},
			}

			by, err = json.Marshal(t)
			if err != nil {
				b.Fail(err)
			}

			rr, err = c.Post(s.URL+router.SignInPath, "application/json", bytes.NewReader(by))
			if err != nil {
				b.Fail(err)
			}

			v := signinmodel.GETResponseSignInModel{}
			if err := json.NewDecoder(rr.Body).Decode(&v); err != nil {
				b.Fail(err)
			}

			b.Assert(v.Service.Ok).Equal(true, "Login is failed")

			rr, err = c.Get(s.URL+router.CheckAuthPath)
			if err != nil {
				b.Fail(err)
			}

			k := checkauthmodel.GETResponseCheckAuthModel{}
			if err := json.NewDecoder(rr.Body).Decode(&k); err != nil {
				b.Fail(err)
			}

			b.Assert(k.Service.Ok).Equal(true, "CheckAuth is failed")
		})
	})
}
