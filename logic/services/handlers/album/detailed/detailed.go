package detailed

// import (
// 	"NewPhotoWeb/logic/proto"
// 	"NewPhotoWeb/logic/services/models"
// 	_ "NewPhotoWeb/utils"
// 	"context"
// 	"encoding/base64"
// 	"encoding/json"
// 	"net/http"

// 	. "NewPhotoWeb/config"
// )

// type IEqualAlbumPage interface {
// 	GetHandler(http.ResponseWriter, *http.Request)
// 	//DeleteHandler(http.ResponseWriter, *http.Request)
// 	//PostHandler(http.ResponseWriter, *http.Request)
// 	//PutHandler(http.ResponseWriter, *http.Request)
// }

// type equalalbumpage struct{}

// func (a *equalalbumpage) GetHandler(w http.ResponseWriter, r *http.Request) {
// 	session, err := Storage.Get(r, "sessionid")

// 	if err != nil {
// 		Logger.Warnln(err.Error())
// 	}

// 	resp := new(models.GetRespEqualAlbumModel)

// 	values, okValue := r.URL.Query()["name"]

// 	userid, okUserid := session.Values["userid"].(string)

// 	if okValue && okUserid {

// 		resp.Service.Ok = true
		
// 		stream, err := NPC.AllPhotosAlbum(context.Background(), &proto.AllPhotosAlbumRequest{Userid: userid, Name: values[0]})
// 		if err != nil {
// 			Logger.ClientError()
// 			return
// 		}
// 		for {
// 			recv, err := stream.Recv()
// 			if err != nil{
// 				break
// 			}

// 			resp.Result = append(resp.Result, struct{
// 				Photo string "json:\"photo\""
// 				Thumbnail string "json:\"thumbnail\""
// 			}{
// 				base64.StdEncoding.EncodeToString(recv.GetPhoto()),
// 				base64.StdEncoding.EncodeToString(recv.GetThumbnail()),
// 			})
// 		}
// 		err = stream.CloseSend()
// 		if err != nil {
// 			Logger.Fatalln(err.Error())
// 		}
		
// 	}else{
// 		resp.Service.Ok = false
// 	}
// 	json.NewEncoder(w).Encode(resp)
// }

// func (a *equalalbumpage) PostHandler(w http.ResponseWriter, r *http.Request) {

// 	r.ParseForm()
// 	albumname := r.FormValue("albumname")

// 	session, err := Storage.Get(r, "sessionid")
// 	if err != nil {
// 		Logger.Warnln(err.Error())
// 	}

// 	var tempvar utils.AlbumPageVars

// 	response, err := NPC.CreateAlbum(context.Background(), &proto.CreateAlbumRequest{Userid: session.Values["userid"].(string), Name: albumname})
// 	if err != nil {
// 		Logger.ClientError()
// 		return
// 	}
// 	switch response.Error {
// 	case "OK":
// 		tempvar.Messages = append(tempvar.Messages, utils.Messages{Type: "Success", Body: fmt.Sprintf("Album %s was successfully created", albumname)})
// 	default:
// 		tempvar.Messages = append(tempvar.Messages, utils.Messages{Type: "Error", Body: fmt.Sprintf("Something went wrong creating %s album", albumname)})
// 	}
// 	http.Redirect(w, r, "/albums", 301)
// }

// func (a *equalalbumpage) DeleteHandler(w http.ResponseWriter, r *http.Request) {
// 	//Post handler for account page ...

// 	r.ParseForm()
// 	albumname := r.FormValue("albumname")

// 	session, err := Storage.Get(r, "sessionid")
// 	if err != nil {
// 		Logger.Warn(err.Error())
// 	}

// 	var tempvar utils.AlbumPageVars

// 	result, err := NPC.DeleteAlbum(context.Background(), &proto.DeleteAlbumRequest{Userid: session.Values["userid"].(string), Name: albumname})
// 	if err != nil {
// 		Logger.ClientError()
// 		return
// 	}

// 	switch result.Error {
// 	case "OK":
// 		tempvar.Messages = append(tempvar.Messages, utils.Messages{Type: "Success", Body: fmt.Sprintf("Album %s was successfully deleted", albumname)})
// 	default:
// 		tempvar.Messages = append(tempvar.Messages, utils.Messages{Type: "Error", Body: fmt.Sprintf("Something went wrong deleting %s album", albumname)})
// 	}
// 	http.Redirect(w, r, "/albums", 301)
// }

// func (a *equalalbumpage) PutHandler(w http.ResponseWriter, r *http.Request) {

// 	reader, err := r.MultipartReader()

// 	if err != nil {
// 		Logger.Fatalln(err.Error())
// 	}

// 	session, err := Storage.Get(r, "sessionid")
// 	if err != nil {
// 		Logger.Warnln(err.Error())
// 	}

// 	stream, err := NPC.UploadPhotoToAlbum(context.Background())

// 	if err != nil {
// 		Logger.Fatalln(err.Error())
// 	}

// 	var albumname string

// 	for {
// 		res, err := reader.NextPart()

// 		if err == io.EOF {
// 			stream.CloseSend()
// 			break
// 		}

// 		switch res.FormName() {
// 		case "album":
// 			var albumsrc bytes.Buffer
// 			_, err := io.Copy(&albumsrc, res)
// 			if err != nil {
// 				Logger.Errorln(err.Error())
// 				continue
// 			}
// 			albumname = albumsrc.String()
// 		case "file":
// 			var file bytes.Buffer
// 			size, err := io.Copy(&file, res)
// 			if err != nil {
// 				Logger.Errorln(err.Error())
// 				continue
// 			}
// 			err = stream.Send(&proto.UploadPhotoToAlbumRequest{Userid: session.Values["userid"].(string), Photo: file.Bytes(), Extension: utils.GetFileExtension(res.FileName()), Size: float64(size), Album: albumname})
// 			if err != nil {
// 				Logger.Errorln(err.Error())
// 				continue
// 			}
// 		case "files":
// 			// if ext := strings.Split(res.FileName(), ".")[1]; ext == "jpeg" || ext == "png" || ext == "jpg" {
// 			// 	var file bytes.Buffer
// 			// 	size, err := io.Copy(&file, res)
// 			// 	if err != nil {
// 			// 		Logger.Errorln(err.Error())
// 			// 		continue
// 			// 	}
// 			// 	err = stream.Send(&proto.UploadPhotoToAlbumRequest{Userid: session.Values["userid"].(string), Photo: file.Bytes(), Extension: utils.GetFileExtension(res.FileName()), Size: float64(size), Album: albumname})
// 			// 	if err != nil {
// 			// 		Logger.Errorln(err.Error(), "HERE")
// 			// 		continue
// 			// 	}
// 			// }
// 		}

// 	}
// 	http.Redirect(w, r, "/albums", 301)
// }

// func NewEqualAlbumPageHandler() IEqualAlbumPage {
// 	return new(equalalbumpage)
// }
