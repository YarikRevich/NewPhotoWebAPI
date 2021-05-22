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

	"github.com/nfnt/resize"
	"google.golang.org/grpc"

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
		session, err := Storage.Get(r, "sessionid")
		if err != nil {
			Logger.Warnln(err.Error())
		}

		grpcResp, err := NPC.GetPhotos(context.Background(), &proto.GetPhotosRequest{Userid: session.Values["userid"].(string)}, grpc.MaxCallRecvMsgSize(32*10e6), grpc.MaxCallSendMsgSize(32*10e6))
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
				Photo     string   "json:\"photo\""
				Thumbnail string   "json:\"thumbnail\""
				Tags      []string "json:\"tags\""
			}{
				base64.StdEncoding.EncodeToString(grpcStreamResp.GetPhoto()),
				base64.StdEncoding.EncodeToString(grpcStreamResp.GetThumbnail()),
				utils.GetCleanTags(grpcStreamResp.GetTags()),
			})
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
	
		session, err := Storage.Get(r, "sessionid")
		if err != nil {
			Logger.Warn(err.Error())
		}

		var req photomodel.POSTRequestPhotoModel
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			Logger.Fatalln(err)
		}

		stream, err := NPC.UploadPhoto(context.Background())
		if err != nil {
			Logger.ClientError()
		}
		for _, value := range req.Data {
			efile, err := base64.StdEncoding.DecodeString(value.File)
			if err != nil {
				Logger.Fatalln(err)
			}

			var img image.Image
			switch value.Extension {
			case "png":
				img, err = png.Decode(bytes.NewReader(efile))
				if err != nil {
					Logger.Fatalln(err)
				}
			case "jpeg":
				img, err = jpeg.Decode(bytes.NewReader(efile))
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
				Userid:    session.Values["userid"].(string),
				Photo:     []byte(base64.StdEncoding.EncodeToString(efile)),
				Thumbnail: thumbnail.Bytes(),
				Extension: value.Extension,
				Size:      value.Size,
			}); err != nil {
				Logger.ClientError()
			}
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
