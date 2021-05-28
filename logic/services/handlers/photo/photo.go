package photo

import (
	"bytes"
	"context"
	"encoding/json"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"

	"github.com/nfnt/resize"

	"NewPhotoWeb/logic/proto"
	photomodel "NewPhotoWeb/logic/services/models/photo"
	"NewPhotoWeb/utils"

	. "NewPhotoWeb/config"
)

type IPhotoPage interface {
	GetHandler() http.Handler
	PostHandler() http.Handler
}

type photo struct{}

func (a *photo) GetHandler() http.Handler {
	//Get handler for photo page ...

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		at, _ := r.Cookie("at")
		lt, _ := r.Cookie("lt")

		grpcResp, err := NPC.GetPhotos(
			context.Background(),
			&proto.GetPhotosRequest{
				AccessToken: at.Value,
				LoginToken:  lt.Value,
			},
		)
		if err != nil {
			Logger.ClientError()
		}
		resp := new(photomodel.GETResponsePhotoModel)

		for {
			grpcStreamResp, err := grpcResp.Recv()
			if err != nil {
				break
			}

			resp.Result = append(resp.Result, struct {
				Thumbnail []byte   "json:\"thumbnail\""
				Tags      []string "json:\"tags\""
			}{
				grpcStreamResp.GetThumbnail(),
				utils.GetCleanTags(grpcStreamResp.GetTags()),
			})
		}

		if err := grpcResp.CloseSend(); err != nil {
			Logger.ClientError()
		}
		resp.Service.Ok = true

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			Logger.Fatalln(err)
		}
	})
}

func (a *photo) PostHandler() http.Handler {
	//Post handler for photo page ...

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		at, _ := r.Cookie("at")
		lt, _ := r.Cookie("lt")

		var req photomodel.POSTRequestPhotoModel
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			Logger.Fatalln(err)
		}

		stream, err := NPC.UploadPhoto(context.Background())
		if err != nil {
			Logger.ClientError()
		}
		for _, value := range req.Data {
			var img image.Image
			switch value.Extension {
			case "png":
				img, err = png.Decode(bytes.NewReader(value.File))
				if err != nil {
					Logger.Fatalln(err)
				}
			case "jpeg":
				img, err = jpeg.Decode(bytes.NewReader(value.File))
				if err != nil {
					Logger.Fatalln(err)
				}
			default:
				continue
			}
			resized := resize.Resize(500, 500, img, resize.Lanczos3)

			var buf bytes.Buffer
			thumbnail := bytes.NewBuffer(buf.Bytes())

			if err = jpeg.Encode(thumbnail, resized, nil); err != nil {
				Logger.Fatalln(err)
			}

			if err = stream.Send(&proto.UploadPhotoRequest{
				AccessToken: at.Value,
				LoginToken:  lt.Value,
				Photo:       value.File,
				Thumbnail:   thumbnail.Bytes(),
				Extension:   value.Extension,
				Size:        value.Size,
			}); err != nil {
				Logger.ClientError()
			}
		}
		if err := stream.CloseSend(); err != nil {
			Logger.Fatalln(err)
		}
		resp := new(photomodel.POSTResponsePhotoModel)
		resp.Service.Ok = true
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			Logger.Fatalln(err)
		}
	})
}

func NewPhotoPageHandler() IPhotoPage {
	return new(photo)
}
