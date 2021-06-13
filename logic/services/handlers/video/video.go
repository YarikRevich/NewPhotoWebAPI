package video

import (
	"NewPhotoWeb/internal"
	"NewPhotoWeb/logic/proto"
	videomodel "NewPhotoWeb/logic/services/models/video"
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"NewPhotoWeb/log"
	"NewPhotoWeb/logic/client"
)

type IVideoPage interface {
	GetHandler() http.Handler
	PostHandler() http.Handler
}

type video struct{}

func (a *video) GetHandler() http.Handler {
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

		grpcResp, err := client.NewPhotoClient.GetVideos(
			context.Background(),
			&proto.GetVideosRequest{
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
		resp := new(videomodel.GETResponseVideoModel)

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
				[]string{},
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

func (a *video) PostHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		at := r.Header["X-At"]
		lt := r.Header["X-Lt"]

		var req videomodel.POSTRequestVideoModel
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Logger.Fatalln(err)
		}

		stream, err := client.NewPhotoClient.UploadVideo(context.Background())
		if err != nil {
			log.Logger.ClientError()
			client.Restart()
		}
		for _, v := range req.Data {
			if err := stream.Send(&proto.UploadVideoRequest{
				AccessToken: at[0],
				LoginToken:  lt[0],
				Video:       v.File,
				Thumbnail:   internal.CreateThumbnailFromVideo(v.File),
				Extension:   v.Extension,
				Size:        v.Size,
			}); err != nil {
				log.Logger.ClientError()
				client.Restart()
			}
		}
		if err := stream.CloseSend(); err != nil {
			log.Logger.ClientError()
			client.Restart()
		}

		resp := new(videomodel.POSTResponseVideoModel)
		resp.Service.Ok = true
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Logger.Fatalln(err)
		}
	})
}

func NewVideoPageHandler() IVideoPage {
	return new(video)
}
