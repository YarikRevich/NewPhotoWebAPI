package video

import (
	"context"
	// "bytes"
	// "context"
	// "encoding/base64"
	"encoding/json"
	// "image"
	// "image/jpeg"
	// "image/png"
	"net/http"

	// "github.com/nfnt/resize"
	// "google.golang.org/grpc"

	// "NewPhotoWeb/logic/proto"
	"NewPhotoWeb/logic/proto"
	videomodel "NewPhotoWeb/logic/services/models/video"

	. "NewPhotoWeb/config"
)

type IVideoPage interface {
	// GetHandler() http.Handler
	PostHandler() http.Handler
}

type video struct{}

// func (a *video) GetHandler() http.Handler {
// 	//Get handler for photo page ...

// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		session, err := Storage.Get(r, "sessionid")
// 		if err != nil {
// 			Logger.Warnln(err.Error())
// 		}

// 	grpcResp, err := NPC.AllPhotos(context.Background(), &proto.AllPhotosRequest{Userid: session.Values["userid"].(string)}, grpc.MaxCallRecvMsgSize(32*10e6), grpc.MaxCallSendMsgSize(32*10e6))
// 	if err != nil {
// 		Logger.ClientError()
// 	}
// 	resp := new(photomodel.GETResponsePhotoModel)

// 	for {
// 		grpcStreamResp, err := grpcResp.Recv()
// 		if err != nil {
// 			break
// 		}

// 		resp.Result = append(resp.Result, struct {
// 			Photo     string   "json:\"photo\""
// 			Thumbnail string   "json:\"thumbnail\""
// 			Tags      []string "json:\"tags\""
// 		}{
// 			base64.StdEncoding.EncodeToString(grpcStreamResp.GetPhoto()),
// 			base64.StdEncoding.EncodeToString(grpcStreamResp.GetThumbnail()),
// 			utils.GetCleanTags(grpcStreamResp.GetTags()),
// 		})
// 	}
// 	resp.Service.Ok = true
// 	if err := json.NewEncoder(w).Encode(resp); err != nil {
// 		Logger.Fatalln(err)
// 	}
// })
// }

func (a *video) PostHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := Storage.Get(r, "sessionid")
		if err != nil {
			Logger.Fatalln(err)
		}

		var req videomodel.POSTRequestVideoModel
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			Logger.Fatalln(err)
		}

		stream, err := NPC.UploadVideo(context.Background())
		if err != nil {
			Logger.ClientError()
		}
		for _, v := range req.Data {
			if err := stream.Send(&proto.UploadVideoRequest{
				Userid:    session.Values["userid"].(string),
				Video:     []byte(v.File),
				Extension: v.Extension,
				Size:      v.Size,
			}); err != nil {
				Logger.ClientError()
			}
		}

		resp := new(videomodel.POSTResponseVideoModel)
		resp.Service.Ok = true
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			Logger.Fatalln(err)
		}
	})
}

func NewVideoPageHandler() IVideoPage {
	return new(video)
}
