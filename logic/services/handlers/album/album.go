package handlers

import (
	"NewPhotoWeb/logic/proto"
	albummodel "NewPhotoWeb/logic/services/models/album"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	. "NewPhotoWeb/config"
)

type IAlbumPage interface {
	GetHandler() http.Handler
	PostHandler() http.Handler
	DeleteHandler() http.Handler
}

type album struct{}

func (a *album) GetHandler() http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		at, _ := r.Cookie("at")
		lt, _ := r.Cookie("lt")

		grpcResp, err := NPC.GetAlbums(
			context.Background(),
			&proto.GetAlbumsRequest{
				AccessToken: at.Value,
				LoginToken:  lt.Value,
			},
		)
		if err != nil {
			Logger.ClientError()
		}

		var resp albummodel.GETResponseAlbumModel

		for {
			grpcStreamResp, err := grpcResp.Recv()
			if err != nil {
				break
			}
			resp.Result = append(resp.Result, struct {
				Name                 string `json:"name"`
				LatestPhoto          []byte `json:"latestphoto"`
				LatestPhotoThumbnail []byte `json:"latestphotothumbnail"`
			}{
				Name:                 grpcStreamResp.GetName(),
				LatestPhoto:          grpcStreamResp.GetLatestPhoto(),
				LatestPhotoThumbnail: grpcStreamResp.GetLatestPhotoThumbnail(),
			})
		}

		resp.Service.Ok = true
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			Logger.Fatalln(err)
		}
	})
}

func (a *album) PostHandler() http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req albummodel.POSTRequestAlbumModel
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			Logger.Fatalln(err)
		}

		at, _ := r.Cookie("at")
		lt, _ := r.Cookie("lt")

		resp := new(albummodel.POSTResponseAlbumModel)

		grpcResp, err := NPC.CreateAlbum(
			context.Background(),
			&proto.CreateAlbumRequest{
				AccessToken: at.Value,
				LoginToken:  lt.Value,
				Name:        req.Data.Name,
			},
		)
		if err != nil {
			Logger.Fatalln(err)
		}
		if grpcResp.GetOk() {
			resp.Service.Message = fmt.Sprintf("Something went wrong creating %s album", req.Data.Name)
		}

		resp.Service.Ok = true
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			Logger.Fatalln(err)
		}
	})
}

func (a *album) DeleteHandler() http.Handler {
	//Post handler for account page ...

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		values, ok := r.URL.Query()["name"]
		if !ok {
			Logger.Fatalln("Album name was not passed")
		}

		at, _ := r.Cookie("at")
		lt, _ := r.Cookie("lt")

		resp := new(albummodel.DELETEResponseAlbumModel)

		grpcResp, err := NPC.DeleteAlbum(
			context.Background(),
			&proto.DeleteAlbumRequest{
				AccessToken: at.Value,
				LoginToken:  lt.Value,
				Name:        values[0],
			},
		)
		if err != nil {
			Logger.Fatalln(err.Error())
		}

		if grpcResp.GetOk() {
			resp.Service.Ok = true
		}

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			Logger.Fatalln(err)
		}
	})
}

// func (a *album) PutHandler() http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		var req albummodel.PUTRequestAlbumModel
// 		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 			Logger.Fatalln(err)
// 		}

// 		at, _ := r.Cookie("at")
// 		lt, _ := r.Cookie("lt")

// 		grpcResp, err := NPC.UploadPhotoToAlbum(context.Background())
// 		if err != nil {
// 			Logger.Fatalln(err)
// 		}

// 		resp := new(albummodel.PUTResponseAlbumModel)

// 		for _, value := range req.Result.Data {

// 			var img image.Image

// 			if value.Extension == "png" {
// 				img, err = png.Decode(bytes.NewReader(value.File))
// 				if err != nil {
// 					Logger.Fatalln(err)
// 				}
// 			} else {
// 				img, err = jpeg.Decode(bytes.NewReader(value.File))
// 				if err != nil {
// 					Logger.Fatalln(err)
// 				}
// 			}

// 			resized := resize.Resize(1280, 800, img, resize.Lanczos3)

// 			var buf bytes.Buffer
// 			thumbnail := bytes.NewBuffer(buf.Bytes())

// 			err = jpeg.Encode(thumbnail, resized, nil)
// 			if err != nil {
// 				Logger.Fatalln(err)
// 			}

// 			if err := grpcResp.Send(
// 				&proto.UploadPhotoToAlbumRequest{
// 					AccessToken: at.Value,
// 					LoginToken:  lt.Value,
// 					Photo:       value.File,
// 					Thumbnail:   thumbnail.Bytes(),
// 					Extension:   value.Extension,
// 					Size:        value.Size,
// 					Album:       req.Result.Name,
// 				},
// 			); err != nil {
// 				Logger.Fatalln(err)
// 			}
// 		}
// 		err = grpcResp.CloseSend()
// 		if err != nil {
// 			Logger.Fatalln(err)
// 		}
// 		resp.Service.Ok = true

// 		if err := json.NewEncoder(w).Encode(resp); err != nil {
// 			Logger.Fatalln(err)
// 		}
// 	})
// }

func NewAlbumHandler() IAlbumPage {
	return new(album)
}
