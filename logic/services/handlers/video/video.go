package video

import (
	"NewPhotoWeb/logic/proto"
	videomodel "NewPhotoWeb/logic/services/models/video"
	"context"
	"encoding/json"
	"net/http"


	"NewPhotoWeb/log"
	"NewPhotoWeb/logic/client"
)

type IVideoPage interface {
	PostHandler() http.Handler
}

type video struct{}

func (a *video) PostHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		at, _ := r.Cookie("at")
		lt, _ := r.Cookie("lt")

		var req videomodel.POSTRequestVideoModel
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Logger.Fatalln(err)
		}

		stream, err := client.NewPhotoClient.UploadVideo(context.Background())
		if err != nil {
			log.Logger.ClientError(); client.Restart()
		}
		for _, v := range req.Data {
			if err := stream.Send(&proto.UploadVideoRequest{
				AccessToken: at.Value,
				LoginToken:  lt.Value,
				Video:       v.File,
				Extension:   v.Extension,
				Size:        v.Size,
			}); err != nil {
				log.Logger.ClientError(); client.Restart()
			}
		}
		if err := stream.CloseSend(); err != nil {
			log.Logger.ClientError(); client.Restart()
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
