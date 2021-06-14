package photo

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	"strconv"

	"github.com/nfnt/resize"

	"NewPhotoWeb/logic/proto"
	photomodel "NewPhotoWeb/logic/services/models/photo"

	"NewPhotoWeb/log"
	"NewPhotoWeb/logic/client"
)

type IPhotoPage interface {
	GetHandler() http.Handler
	PostHandler() http.Handler
}

type photo struct{}

func (a *photo) GetHandler() http.Handler {
	//Get handler for photo page ...

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		at := r.Header["X-At"]
		lt := r.Header["X-Lt"]

		offset, ok := r.URL.Query()["offset"]
		if !ok {
			log.Logger.Fatalln("Offset is empty!")
		}
		offsetNum, err := strconv.Atoi(offset[0])
		if err != nil {
			log.Logger.Fatalln(err)
		}
		page, ok := r.URL.Query()["page"]
		if !ok {
			log.Logger.Fatalln("Page is empty!")
		}
		pageNum, err := strconv.Atoi(page[0])
		if err != nil {
			log.Logger.Fatalln(err)
		}

		grpcResp, err := client.NewPhotoClient.GetPhotos(
			context.Background(),
			&proto.GetPhotosRequest{
				AccessToken: at[0],
				LoginToken:  lt[0],
				Offset:      int64(offsetNum),
				Page:        int64(pageNum),
			},
		)
		if err != nil {
			log.Logger.ClientError()
			client.Restart()
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
				grpcStreamResp.GetTags(),
			})
		}

		if err := grpcResp.CloseSend(); err != nil {
			log.Logger.ClientError()
			client.Restart()
		}

		resp.Service.Ok = true

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Logger.Fatalln(err)
		}
	})
}

func (a *photo) PostHandler() http.Handler {
	//Post handler for photo page ...

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		at := r.Header["X-At"]
		lt := r.Header["X-Lt"]

		var req photomodel.POSTRequestPhotoModel
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Logger.Fatalln(err)
		}

		stream, err := client.NewPhotoClient.UploadPhoto(context.Background())
		if err != nil {
			log.Logger.ClientError()
			client.Restart()
		}

		for _, value := range req.Data {
			var img image.Image
			var file = value.File
			if b, err := base64.StdEncoding.DecodeString(string(value.File)); err == nil {
				file = b
			}

			switch value.Extension {
			case "png":
				img, err = png.Decode(bytes.NewReader(file))
				if err != nil {
					log.Logger.Fatalln(err)
				}
			case "jpeg":
				img, err = jpeg.Decode(bytes.NewReader(file))
				if err != nil {
					log.Logger.Fatalln(err)
				}
			default:
				continue
			}

			resized := resize.Resize(500, 500, img, resize.Lanczos3)

			var buf bytes.Buffer
			thumbnail := bytes.NewBuffer(buf.Bytes())

			switch value.Extension {
			case "png":
				if err = png.Encode(thumbnail, resized); err != nil {
					log.Logger.Fatalln(err)
				}
			case "jpeg":
				if err = jpeg.Encode(thumbnail, resized, nil); err != nil {
					log.Logger.Fatalln(err)
				}
			}

			if err = stream.Send(&proto.UploadPhotoRequest{
				AccessToken: at[0],
				LoginToken:  lt[0],
				Photo:       value.File,
				Thumbnail:   thumbnail.Bytes(),
				Extension:   value.Extension,
				Size:        value.Size,
			}); err != nil {
				log.Logger.ClientError()
				client.Restart()
			}
		}
		if err := stream.CloseSend(); err != nil {
			log.Logger.ClientError()
			client.Restart()
		}
		resp := new(photomodel.POSTResponsePhotoModel)
		resp.Service.Ok = true
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Logger.Fatalln(err)
		}
	})
}

func NewPhotoPageHandler() IPhotoPage {
	return new(photo)
}
