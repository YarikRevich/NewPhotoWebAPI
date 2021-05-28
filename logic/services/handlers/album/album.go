package handlers

import (
	"NewPhotoWeb/logic/proto"
	albummodel "NewPhotoWeb/logic/services/models/album"
	"context"
	"encoding/json"
	"fmt"
	"net/http"


	"NewPhotoWeb/log"
	"NewPhotoWeb/logic/client"
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

		grpcResp, err := client.NewPhotoClient.GetAlbums(
			context.Background(),
			&proto.GetAlbumsRequest{
				AccessToken: at.Value,
				LoginToken:  lt.Value,
			},
		)
		if err != nil {
			log.Logger.ClientError(); client.Restart()
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
			log.Logger.Fatalln(err)
		}
	})
}

func (a *album) PostHandler() http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req albummodel.POSTRequestAlbumModel
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Logger.Fatalln(err)
		}

		at, _ := r.Cookie("at")
		lt, _ := r.Cookie("lt")

		resp := new(albummodel.POSTResponseAlbumModel)

		grpcResp, err := client.NewPhotoClient.CreateAlbum(
			context.Background(),
			&proto.CreateAlbumRequest{
				AccessToken: at.Value,
				LoginToken:  lt.Value,
				Name:        req.Data.Name,
			},
		)
		if err != nil {
			log.Logger.ClientError(); client.Restart()
		}
		if grpcResp.GetOk() {
			resp.Service.Message = fmt.Sprintf("Something went wrong creating %s album", req.Data.Name)
		}

		resp.Service.Ok = true
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Logger.Fatalln(err)
		}
	})
}

func (a *album) DeleteHandler() http.Handler {
	//Post handler for account page ...

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		values, ok := r.URL.Query()["name"]
		if !ok {
			log.Logger.Fatalln("Name of album is not passed")
		}

		at, _ := r.Cookie("at")
		lt, _ := r.Cookie("lt")

		resp := new(albummodel.DELETEResponseAlbumModel)

		grpcResp, err := client.NewPhotoClient.DeleteAlbum(
			context.Background(),
			&proto.DeleteAlbumRequest{
				AccessToken: at.Value,
				LoginToken:  lt.Value,
				Name:        values[0],
			},
		)
		if err != nil {
			log.Logger.ClientError(); client.Restart()
		}

		if grpcResp.GetOk() {
			resp.Service.Ok = true
		}

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Logger.Fatalln(err)
		}
	})
}

func NewAlbumHandler() IAlbumPage {
	return new(album)
}
