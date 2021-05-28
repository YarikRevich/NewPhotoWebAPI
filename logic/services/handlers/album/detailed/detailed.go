package detailed

import (
	"context"
	"encoding/json"
	"net/http"

	. "NewPhotoWeb/config"
	"NewPhotoWeb/logic/proto"
	detailedalbummodel "NewPhotoWeb/logic/services/models/album/detailed"
)

type IDetailedAlbumPage interface {
	GetHandler() http.Handler
	DeleteHandler() http.Handler
	PutHandler() http.Handler
}

type detailedalbum struct{}

func (a *detailedalbum) GetHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		values, ok := r.URL.Query()["name"]
		if !ok {
			Logger.Fatalln("Album name is empty!")
		}

		at, _ := r.Cookie("at")
		lt, _ := r.Cookie("lt")

		resp := new(detailedalbummodel.GETResponseEqualAlbumModel)

		grpcStreamPhotosResp, err := NPC.GetPhotosFromAlbum(
			context.Background(),
			&proto.GetPhotosFromAlbumRequest{
				AccessToken: at.Value,
				LoginToken:  lt.Value,
				Name:        values[0],
			},
		)
		if err != nil {
			Logger.ClientError()
		}

		for {
			recv, err := grpcStreamPhotosResp.Recv()
			if err != nil {
				break
			}
			resp.Result.Photos = append(resp.Result.Photos, struct {
				Photo     []byte "json:\"photo\""
				Thumbnail []byte "json:\"thumbnail\""
				Extension string "json:\"extension\""
			}{
				recv.GetPhoto(),
				recv.GetThumbnail(),
				recv.GetExtension(),
			})
		}
		if err = grpcStreamPhotosResp.CloseSend(); err != nil {
			Logger.Fatalln(err.Error())
		}

		grpcStreamVideosResp, err := NPC.GetVideosFromAlbum(
			context.Background(),
			&proto.GetVideosFromAlbumRequest{
				AccessToken: at.Value,
				LoginToken:  lt.Value,
				Name:        values[0],
			},
		)
		if err != nil {
			Logger.ClientError()
		}

		for {
			recv, err := grpcStreamVideosResp.Recv()
			if err != nil {
				break
			}

			resp.Result.Videos = append(resp.Result.Videos, struct {
				Video     []byte "json:\"video\""
				Extension string "json:\"extension\""
			}{
				recv.GetVideo(),
				recv.GetExtension(),
			})
		}
		if err = grpcStreamVideosResp.CloseSend(); err != nil {
			Logger.Fatalln(err.Error())
		}

		resp.Result.Name = values[0]
		resp.Service.Ok = true

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			Logger.Fatalln(err)
		}
	})
}

func (a *detailedalbum) DeleteHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		at, _ := r.Cookie("at")
		lt, _ := r.Cookie("lt")

		var req detailedalbummodel.DELETERequestEqualAlbumModel
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			Logger.Fatalln(err)
		}

		grpcRespPhotoStream, err := NPC.DeletePhotoFromAlbum(context.Background())
		if err != nil {
			Logger.ClientError()
		}

		grpcRespVideoStream, err := NPC.DeleteVideoFromAlbum(context.Background())
		if err != nil {
			Logger.ClientError()
		}

		for _, v := range req.Data.Photos {
			if err := grpcRespPhotoStream.Send(
				&proto.DeletePhotoFromAlbumRequest{
					AccessToken: at.Value,
					LoginToken:  lt.Value,
					Photo:       v,
					Album:       req.Data.Name,
				},
			); err != nil {
				Logger.Fatalln(err)
			}
		}

		for _, v := range req.Data.Videos {
			if err := grpcRespVideoStream.Send(
				&proto.DeleteVideoFromAlbumRequest{
					AccessToken: at.Value,
					LoginToken:  lt.Value,
					Video:       v,
					Album:       req.Data.Name,
				},
			); err != nil {
				Logger.Fatalln(err)
			}
		}

		grpcRespPhoto, err := grpcRespPhotoStream.CloseAndRecv()
		if err != nil {
			Logger.ClientError()
		}

		grpcRespVideo, err := grpcRespVideoStream.CloseAndRecv()
		if err != nil {
			Logger.ClientError()
		}

		var resp detailedalbummodel.DELETEResponseEqualAlbumModel
		if grpcRespPhoto.GetOk() && grpcRespVideo.GetOk() {
			resp.Service.Ok = true
		}
		if err := json.NewEncoder(w).Encode(&resp); err != nil {
			Logger.Fatalln(err)
		}
	})
}

func (a *detailedalbum) PutHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		at, _ := r.Cookie("at")
		lt, _ := r.Cookie("lt")

		var req detailedalbummodel.PUTRequestEqualAlbumModel
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			Logger.Fatalln(err)
		}

		streamImage, err := NPC.UploadPhotoToAlbum(context.Background())
		if err != nil {
			Logger.ClientError()
		}

		streamVideo, err := NPC.UploadVideoToAlbum(context.Background())
		if err != nil {
			Logger.ClientError()
		}

		for _, v := range req.Data.Photos {
			if err = streamImage.Send(
				&proto.UploadPhotoToAlbumRequest{
					AccessToken: at.Value,
					LoginToken:  lt.Value,
					Photo:       v.File,
					Extension:   v.Extension,
					Size:        float64(v.Size),
					Album:       req.Data.Name,
				},
			); err != nil {
				Logger.Fatalln(err.Error())
			}
		}

		for _, v := range req.Data.Videos {
			if err = streamVideo.Send(
				&proto.UploadVideoToAlbumRequest{
					AccessToken: at.Value,
					LoginToken:  lt.Value,
					Video:       v.File,
					Extension:   v.Extension,
					Size:        float64(v.Size),
					Album:       req.Data.Name,
				},
			); err != nil {
				Logger.Fatalln(err.Error())
			}
		}

		grpcRespPhotos, err := streamImage.CloseAndRecv()
		if err != nil {
			Logger.Fatalln(err)
		}
		grpcRespVideos, err := streamVideo.CloseAndRecv()
		if err != nil {
			Logger.Fatalln(err)
		}
		var resp detailedalbummodel.PUTResponseEqualAlbumModel
		if grpcRespPhotos.GetOk() && grpcRespVideos.GetOk() {
			resp.Service.Ok = true
		}

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			Logger.Fatalln(err)
		}
	})
}

func NewDetailedAlbumPageHandler() IDetailedAlbumPage {
	return new(detailedalbum)
}
