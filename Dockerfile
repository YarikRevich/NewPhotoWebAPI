FROM golang:latest

ENV serverAddr="newphoto_server:8082"

WORKDIR /go/src/NewPhotoWeb

COPY . .

RUN go build main.go

ENTRYPOINT ./main
