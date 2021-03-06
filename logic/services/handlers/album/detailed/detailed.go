package detailed

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

	"NewPhotoWeb/internal"
	"NewPhotoWeb/log"
	"NewPhotoWeb/logic/client"
	"NewPhotoWeb/logic/proto"
	detailedalbummodel "NewPhotoWeb/logic/services/models/album/detailed"
	"github.com/nfnt/resize"
)

type IDetailedAlbumPage interface {
	GetHandler() http.Handler
	DeleteHandler() http.Handler
	PutHandler() http.Handler
}

type detailedalbum struct{}

func (a *detailedalbum) GetHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name, ok := r.URL.Query()["name"]
		if !ok {
			log.Logger.Fatalln("Album name is empty!")
		}
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

		at := r.Header["X-At"]
		lt := r.Header["X-Lt"]

		resp := new(detailedalbummodel.GETResponseEqualAlbumModel)

		grpcStreamPhotosResp, err := client.NewPhotoClient.GetPhotosFromAlbum(
			context.Background(),
			&proto.GetPhotosFromAlbumRequest{
				AccessToken: at[0],
				LoginToken:  lt[0],
				Name:        name[0],
				Offset:      int64(offsetNum),
				Page:        int64(pageNum),
			},
		)
		if err != nil {
			log.Logger.ClientError()
			client.Restart()
		}

		for {
			recv, err := grpcStreamPhotosResp.Recv()
			if err != nil {
				break
			}
			resp.Result.Photos = append(resp.Result.Photos, struct {
				Thumbnail []byte   "json:\"thumbnail\""
				Tags      []string "json:\"tags\""
			}{
				recv.GetThumbnail(),
				recv.GetTags(),
			})
		}
		if err = grpcStreamPhotosResp.CloseSend(); err != nil {
			log.Logger.ClientError()
			client.Restart()
		}

		grpcStreamVideosResp, err := client.NewPhotoClient.GetVideosFromAlbum(
			context.Background(),
			&proto.GetVideosFromAlbumRequest{
				AccessToken: at[0],
				LoginToken:  lt[0],
				Name:        name[0],
				Offset:      int64(offsetNum),
				Page:        int64(pageNum),
			},
		)
		if err != nil {
			log.Logger.ClientError()
			client.Restart()
		}

		for {
			recv, err := grpcStreamVideosResp.Recv()
			if err != nil {
				break
			}

			resp.Result.Videos = append(resp.Result.Videos, struct {
				Thumbnail []byte   "json:\"thumbnail\""
				Tags      []string "json:\"tags\""
			}{
				recv.GetThumbnail(),
				recv.GetTags(),
			})
		}
		if err = grpcStreamVideosResp.CloseSend(); err != nil {
			log.Logger.ClientError()
			client.Restart()
		}

		resp.Result.Name = name[0]
		resp.Service.Ok = true

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Logger.Fatalln(err)
		}
	})
}

func (a *detailedalbum) DeleteHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		at := r.Header["X-At"]
		lt := r.Header["X-Lt"]

		var req detailedalbummodel.DELETERequestEqualAlbumModel
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Logger.Fatalln(err)
		}

		grpcRespPhotoStream, err := client.NewPhotoClient.DeletePhotoFromAlbum(context.Background())
		if err != nil {
			log.Logger.ClientError()
			client.Restart()
		}

		grpcRespVideoStream, err := client.NewPhotoClient.DeleteVideoFromAlbum(context.Background())
		if err != nil {
			log.Logger.ClientError()
			client.Restart()
		}

		for _, v := range req.Data.Photos {
			if err := grpcRespPhotoStream.Send(
				&proto.DeletePhotoFromAlbumRequest{
					AccessToken: at[0],
					LoginToken:  lt[0],
					Photo:       v,
					Album:       req.Data.Name,
				},
			); err != nil {
				log.Logger.ClientError()
				client.Restart()
			}
		}

		for _, v := range req.Data.Videos {
			if err := grpcRespVideoStream.Send(
				&proto.DeleteVideoFromAlbumRequest{
					AccessToken: at[0],
					LoginToken:  lt[0],
					Video:       v,
					Album:       req.Data.Name,
				},
			); err != nil {
				log.Logger.ClientError()
				client.Restart()
			}
		}

		grpcRespPhoto, err := grpcRespPhotoStream.CloseAndRecv()
		if err != nil {
			log.Logger.ClientError()
			client.Restart()
		}

		grpcRespVideo, err := grpcRespVideoStream.CloseAndRecv()
		if err != nil {
			log.Logger.ClientError()
			client.Restart()
		}

		var resp detailedalbummodel.DELETEResponseEqualAlbumModel
		if grpcRespPhoto.GetOk() && grpcRespVideo.GetOk() {
			resp.Service.Ok = true
		}
		if err := json.NewEncoder(w).Encode(&resp); err != nil {
			log.Logger.Fatalln(err)
		}
	})
}

func (a *detailedalbum) PutHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		at := r.Header["X-At"]
		lt := r.Header["X-Lt"]

		var req detailedalbummodel.PUTRequestEqualAlbumModel
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Logger.Fatalln(err)
		}

		streamImage, err := client.NewPhotoClient.UploadPhotoToAlbum(context.Background())
		if err != nil {
			log.Logger.ClientError()
			client.Restart()
		}

		streamVideo, err := client.NewPhotoClient.UploadVideoToAlbum(context.Background())
		if err != nil {
			log.Logger.ClientError()
			client.Restart()
		}

		for _, v := range req.Data.Photos {
			var img image.Image
			var file = v.File
			if b, err := base64.StdEncoding.DecodeString(string(v.File)); err == nil {
				file = b
			}

			switch v.Extension {
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

			switch v.Extension {
			case "png":
				if err = png.Encode(thumbnail, resized); err != nil {
					log.Logger.Fatalln(err)
				}
			case "jpeg":
				if err = jpeg.Encode(thumbnail, resized, nil); err != nil {
					log.Logger.Fatalln(err)
				}
			}

			if err = streamImage.Send(
				&proto.UploadPhotoToAlbumRequest{
					AccessToken: at[0],
					LoginToken:  lt[0],
					Photo:       v.File,
					Thumbnail:   thumbnail.Bytes(),
					Extension:   v.Extension,
					Size:        float64(v.Size),
					Album:       req.Data.Name,
				},
			); err != nil {
				log.Logger.ClientError()
				client.Restart()
			}
		}

		for _, v := range req.Data.Videos {
			if err = streamVideo.Send(
				&proto.UploadVideoToAlbumRequest{
					AccessToken: at[0],
					LoginToken:  lt[0],
					Video:       v.File,
					Thumbnail:   internal.CreateThumbnailFromVideo(v.File),
					Extension:   v.Extension,
					Size:        float64(v.Size),
					Album:       req.Data.Name,
				},
			); err != nil {
				log.Logger.ClientError()
				client.Restart()
			}
		}

		grpcRespPhotos, err := streamImage.CloseAndRecv()
		if err != nil {
			log.Logger.ClientError()
			client.Restart()
		}
		grpcRespVideos, err := streamVideo.CloseAndRecv()
		if err != nil {
			log.Logger.ClientError()
			client.Restart()
		}
		var resp detailedalbummodel.PUTResponseEqualAlbumModel
		if grpcRespPhotos.GetOk() && grpcRespVideos.GetOk() {
			resp.Service.Ok = true
		}

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Logger.Fatalln(err)
		}
	})
}

func NewDetailedAlbumPageHandler() IDetailedAlbumPage {
	return new(detailedalbum)
}
