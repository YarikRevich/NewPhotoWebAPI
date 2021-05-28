package tests

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	accountmodel "NewPhotoWeb/logic/services/models/account"
	signinmodel "NewPhotoWeb/logic/services/models/auth/sign_in"
	signupmodel "NewPhotoWeb/logic/services/models/auth/sign_up"
	photomodel "NewPhotoWeb/logic/services/models/photo"
	photodetailedmodel "NewPhotoWeb/logic/services/models/photo/detailed"
	"NewPhotoWeb/logic/services/router"

	"github.com/franela/goblin"
)

func TestPhotoHandler(t *testing.T) {
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

	i, err := os.ReadFile("test.png")
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

	b.Describe("Add photo", func() {
		b.It("Should do that correctly", func() {
			r := photomodel.POSTRequestPhotoModel{}
			r.Data = append(r.Data, struct {
				File      []byte  "json:\"file\""
				Name      string  "json:\"name\""
				Size      float64 "json:\"size\""
				Extension string  "json:\"extension\""
			}{
				File:      i,
				Name:      y,
				Size:      1,
				Extension: "png",
			})

			by, err := json.Marshal(r)
			if err != nil {
				b.Fail(err)
			}

			rr, err := c.Post(s.URL+router.PhotosPath, "application/json", bytes.NewReader(by))
			if err != nil {
				b.Fail(err)
			}
			a := photomodel.POSTResponsePhotoModel{}
			if err := json.NewDecoder(rr.Body).Decode(&a); err != nil {
				b.Fail(err)
			}

			b.Assert(a.Service.Ok).Equal(true, "Photo adding failed")
		})
	})

	time.Sleep(4 * time.Second)

	b.Describe("Get photo", func() {
		b.It("Should do that correctly", func() {
			rr, err := c.Get(s.URL + router.PhotosPath)
			if err != nil {
				b.Fail(err)
			}
			a := photomodel.GETResponsePhotoModel{}
			if err := json.NewDecoder(rr.Body).Decode(&a); err != nil {
				b.Fail(err)
			}

			b.Assert(a.Result[0].Thumbnail).IsNotNil("Thumbnail is not gotten")

			time.Sleep(2 * time.Second)

			r := photodetailedmodel.GETRequestDetailedPhotoModel{}
			r.Data.Thumbnail = a.Result[0].Thumbnail

			by, err := json.Marshal(r)
			if err != nil {
				b.Fail(err)
			}

			rr, err = c.Post(s.URL+router.PhotosDetailedPath, "application/json", bytes.NewReader(by))
			if err != nil {
				b.Fail(err)
			}

			q := photodetailedmodel.GETResponseDetailedPhotoModel{}
			if err := json.NewDecoder(rr.Body).Decode(&q); err != nil {
				b.Fail(err)
			}
			b.Assert(q.Result.Photo).Equal(i, "Thumbnail is not gotten")
			b.Assert(q.Service.Ok).IsTrue("Detailed photo getting failed")
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
