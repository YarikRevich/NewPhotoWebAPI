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

	accountmodel "NewPhotoWeb/logic/services/models/account"
	accountavatarmodel "NewPhotoWeb/logic/services/models/account/avatar"
	signinmodel "NewPhotoWeb/logic/services/models/auth/sign_in"
	signupmodel "NewPhotoWeb/logic/services/models/auth/sign_up"
	"NewPhotoWeb/logic/services/router"

	"github.com/franela/goblin"
)

func TestAccountAvatarHandler(t *testing.T) {
	b := goblin.Goblin(t)
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

	y := "Test" + randInt.String()

	b.Describe("Check sign up", func() {
		b.It("Should be signed up", func() {

			q := signupmodel.POSTRequestRegestrationModel{
				Data: struct {
					Login      string "json:\"login\""
					Firstname  string "json:\"firstname\""
					Secondname string "json:\"secondname\""
					Password1  string "json:\"password1\""
					Password2  string "json:\"password2\""
				}{
					Login:      y,
					Firstname:  y,
					Secondname: y,
					Password1:  y,
					Password2:  y,
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
		})
	})

	b.Describe("Check sign in", func() {
		b.It("Should be signed in", func() {
			t := signinmodel.GETRequestSignInModel{
				Data: struct {
					Login    string "json:\"login\""
					Password string "json:\"password\""
				}{
					Login:    y,
					Password: y,
				},
			}

			by, err := json.Marshal(t)
			if err != nil {
				b.Fail(err)
			}

			rr, err := c.Post(s.URL+router.SignInPath, "application/json", bytes.NewReader(by))
			if err != nil {
				b.Fail(err)
			}

			v := signinmodel.GETResponseSignInModel{}
			if err := json.NewDecoder(rr.Body).Decode(&v); err != nil {
				b.Fail(err)
			}

			b.Assert(v.Service.Ok).Equal(true, "Login is failed")
		})
	})

	b.Describe("Sets avatar", func() {
		b.It("Should do that correctly", func() {
			r := accountavatarmodel.POSTRequestAvatarModel{}
			r.Data.Avatar = y
			by, err := json.Marshal(r)
			if err != nil {
				b.Fail(err)
			}

			rr, err := c.Post(s.URL+router.AccountAvatarPath, "application/json", bytes.NewReader(by))
			if err != nil {
				b.Fail(err)
			}

			v := accountavatarmodel.POSTResponseAvatarModel{}
			if err := json.NewDecoder(rr.Body).Decode(&v); err != nil {
				b.Fail(err)
			}

			b.Assert(v.Service.Ok).Equal(true, "Avatar setting is failed")
		})
	})

	b.Describe("Gets avatar", func() {
		b.It("Should do that correctly", func() {
			rr, err := c.Get(s.URL + router.AccountAvatarPath)
			if err != nil {
				b.Fail(err)
			}

			v := accountavatarmodel.GETResponseAvatarModel{}
			if err := json.NewDecoder(rr.Body).Decode(&v); err != nil {
				b.Fail(err)
			}

			b.Assert(v.Result.Avatar).Equal(y, "Avatar is not correct")
			b.Assert(v.Service.Ok).Equal(true, "Avtar setting is failed")
		})
	})

	b.Describe("Delete account", func() {
		b.It("Should do that correctly", func() {
			r, err := http.NewRequest("DELETE", s.URL+router.AccountPath, &bytes.Reader{})
			if err != nil {
				b.Fail(err)
			}
			rr, err := c.Do(r)
			if err != nil {
				b.Fail(err)
			}

			a := accountmodel.DELETEResponseAccountModel{}
			if err := json.NewDecoder(rr.Body).Decode(&a); err != nil {
				b.Fail(err)
			}

			b.Assert(a.Service.Ok).Equal(true, "Account deletion failed")
		})
	})
}
