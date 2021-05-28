package tests

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	accountmodel "NewPhotoWeb/logic/services/models/account"
	albummodel "NewPhotoWeb/logic/services/models/album"
	infodetailedalbummodel "NewPhotoWeb/logic/services/models/album/detailed/info"
	signinmodel "NewPhotoWeb/logic/services/models/auth/sign_in"
	signupmodel "NewPhotoWeb/logic/services/models/auth/sign_up"

	"NewPhotoWeb/logic/services/router"

	"github.com/franela/goblin"
)

func TestAlbumHandler(t *testing.T) {
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

	b.Describe("Create album", func() {
		b.It("Should do that correctly", func() {
			r := albummodel.POSTRequestAlbumModel{}
			r.Data.Name = y

			by, err := json.Marshal(r)
			if err != nil {
				b.Fail(err)
			}

			rr, err := c.Post(s.URL+router.AlbumsPath, "application/json", bytes.NewReader(by))
			if err != nil {
				b.Fail(err)
			}
			a := albummodel.POSTResponseAlbumModel{}
			if err := json.NewDecoder(rr.Body).Decode(&a); err != nil {
				b.Fail(err)
			}

			b.Assert(a.Service.Ok).Equal(true, "Album creation failed")
		})
	})

	time.Sleep(1 * time.Second)

	b.Describe("Get album", func() {
		b.It("Should do that correctly", func() {
			rr, err := c.Get(s.URL + router.AlbumsPath)
			if err != nil {
				b.Fail(err)
			}
			a := albummodel.GETResponseAlbumModel{}
			if err := json.NewDecoder(rr.Body).Decode(&a); err != nil {
				b.Fail(err)
			}
			b.Assert(a.Result[0].Name).Equal(y, "Album name is not correct")
			b.Assert(a.Service.Ok).IsTrue("Album getting failed")
		})
	})

	time.Sleep(2 * time.Second)

	b.Describe("Add photo to album", func() {
		b.It("Should do that correctly", func() {
			r := albummodel.PUTRequestAlbumModel{}
			r.Result.Name = y
			r.Result.Data = append(r.Result.Data, struct {
				File      []byte  "json:\"file\""
				Size      float64 "json:\"size\""
				Extension string  "json:\"extension\""
			}{
				File:      i,
				Size:      1,
				Extension: "png",
			})
			by, err := json.Marshal(r)
			if err != nil {
				b.Fail(err)
			}

			rr, err := http.NewRequest("PUT", s.URL+router.AlbumsPath, bytes.NewReader(by))
			if err != nil {
				b.Fail(err)
			}

			rq, err := c.Do(rr)
			if err != nil {
				b.Fail(err)
			}

			g := albummodel.PUTResponseAlbumModel{}
			if err := json.NewDecoder(rq.Body).Decode(&g); err != nil {
				b.Fail(err)
			}

			b.Assert(g.Service.Ok).IsTrue("Adding to album is not correct")

			time.Sleep(3 * time.Second)

			l, err := c.Get(s.URL + router.AlbumsPath)
			if err != nil {
				b.Fail(err)
			}
			a := albummodel.GETResponseAlbumModel{}
			if err := json.NewDecoder(l.Body).Decode(&a); err != nil {
				b.Fail(err)
			}
			b.Assert(a.Result[0].LatestPhoto).Equal(i, "Album latest photo is not correct")
			b.Assert(a.Service.Ok).IsTrue("Album getting failed")
		})
	})

	b.Describe("Check new album info", func() {
		b.It("Should do that correctly", func() {
			l, err := c.Get(s.URL + router.AlbumsDetailedInfoPath + fmt.Sprintf("?name=%s", y))
			if err != nil {
				b.Fail(err)
			}
			a := infodetailedalbummodel.GETResponseGetAlbumInfoModel{}
			if err := json.NewDecoder(l.Body).Decode(&a); err != nil {
				b.Fail(err)
			}
			b.Assert(a.Result.MediaNum).Equal(int64(1), "A num of media in the album is incorrect")
			b.Assert(a.Service.Ok).IsTrue("Album getting failed")
		})
	})

	b.Describe("Delete album", func() {
		b.It("Should do that correctly", func() {
			rr, err := http.NewRequest("DELETE", s.URL+router.AlbumsPath+fmt.Sprintf("?name=%s", y), nil)
			if err != nil {
				b.Fail(err)
			}

			rq, err := c.Do(rr)
			if err != nil {
				b.Fail(err)
			}

			g := albummodel.DELETEResponseAlbumModel{}
			if err := json.NewDecoder(rq.Body).Decode(&g); err != nil {
				b.Fail(err)
			}

			b.Assert(g.Service.Ok).IsTrue("Album deletion failed")
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
