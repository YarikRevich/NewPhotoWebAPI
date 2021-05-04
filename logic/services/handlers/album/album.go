package handlers

import (
	"NewPhotoWeb/logic/proto"
	albummodel "NewPhotoWeb/logic/services/models/album"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"

	"github.com/nfnt/resize"

	. "NewPhotoWeb/config"
)

type IAlbumPage interface {
	GetHandler() http.Handler
	PostHandler() http.Handler
	PutHandler() http.Handler
	DeleteHandler() http.Handler
}

type album struct{}

func (a *album) GetHandler() http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		session, err := Storage.Get(r, "sessionid")
		if err != nil {
			Logger.Warnln(err.Error())
		}

		grpcResp, err := NPC.GetAllAlbums(context.Background(), &proto.GetAllAlbumsRequest{Userid: session.Values["userid"].(string)})
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
				LatestPhoto          string `json:"latestphoto"`
				LatestPhotoThumbnail string "json:\"latestphotothumbnail\""
			}{
				Name:                 grpcStreamResp.GetName(),
				LatestPhoto:          base64.StdEncoding.EncodeToString(grpcStreamResp.GetLatestPhoto()),
				LatestPhotoThumbnail: base64.StdEncoding.EncodeToString(grpcStreamResp.GetLatestPhotoThumbnail()),
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

		session, err := Storage.Get(r, "sessionid")
		if err != nil {
			Logger.Warnln(err)
		}

		resp := new(albummodel.POSTResponseAlbumModel)

		grpcResp, err := NPC.CreateAlbum(context.Background(), &proto.CreateAlbumRequest{Userid: session.Values["userid"].(string), Name: req.Data.Name})
		if err != nil {
			Logger.Fatalln(err)
		}
		if grpcResp.GetError() != "OK" {
			resp.Service.Message = fmt.Sprintf("Something went wrong creating %s album", req.Data.Name)
		}

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			Logger.Fatalln(err)
		}
	})
}

func (a *album) DeleteHandler() http.Handler {
	//Post handler for account page ...

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var req albummodel.DELETERequestAlbumModel
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			Logger.Fatalln(err)
		}

		session, err := Storage.Get(r, "sessionid")
		if err != nil {
			Logger.Warn(err.Error())
		}

		resp := new(albummodel.DELETEResponseAlbumModel)

		grpcResp, err := NPC.DeleteAlbum(context.Background(), &proto.DeleteAlbumRequest{Userid: session.Values["userid"].(string), Name: req.Data.Name})
		if err != nil {
			Logger.Fatalln(err.Error())
		}

		if grpcResp.GetError() != "OK" {
			resp.Service.Message = fmt.Sprintf("Something went wrong deleting %s album", req.Data.Name)
		}

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			Logger.Fatalln(err)
		}
	})
}

func (a *album) PutHandler() http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var req albummodel.PUPRequestAlbumModel
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil{
			Logger.Fatalln(err)
		}

		session, err := Storage.Get(r, "sessionid")
		if err != nil {
			Logger.Warnln(err.Error())
		}

		grpcResp, err := NPC.UploadPhotoToAlbum(context.Background())
		if err != nil {
			Logger.Fatalln(err)
		}

		resp := new(albummodel.PUPResponseAlbumModel)

		for _, value := range req.Result.Data {
			efile, err := base64.StdEncoding.DecodeString(value.File)
			if err != nil {
				Logger.Fatalln(err)
			}

			var img image.Image

			if value.Extension == "png" {
				img, err = png.Decode(bytes.NewReader(efile))
				if err != nil {
					Logger.Fatalln(err)
				}
			} else {
				img, err = jpeg.Decode(bytes.NewReader(efile))
				if err != nil {
					Logger.Fatalln(err)
				}
			}

			resized := resize.Resize(1280, 800, img, resize.Lanczos3)

			var buf bytes.Buffer
			thumbnail := bytes.NewBuffer(buf.Bytes())

			err = jpeg.Encode(thumbnail, resized, nil)
			if err != nil {
				Logger.Fatalln(err)
			}

			if err := grpcResp.Send(&proto.UploadPhotoToAlbumRequest{Userid: session.Values["userid"].(string), Photo: efile, Thumbnail: thumbnail.Bytes(), Extension: value.Extension, Size: value.Size, Album: req.Result.Name}); err != nil{
				Logger.Fatalln(err)
			}
		}
		err = grpcResp.CloseSend()
		if err != nil {
			Logger.Fatalln(err)
		}
		resp.Service.Ok = true

		if err := json.NewEncoder(w).Encode(resp); err != nil{
			Logger.Fatalln(err)
		}
	})
}

func NewAlbumHandler() IAlbumPage {
	return new(album)
}
