package client

import (
	"log"
	"NewPhotoWeb/logic/proto"
	"fmt"
	"google.golang.org/grpc"
	"os"
	"os/signal"
	"syscall"
)

var (
	conn = NewConnection()
)

func NewConnection() *grpc.ClientConn {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(50 * 10e6), grpc.MaxCallSendMsgSize(50 * 10e6)),
	}

	serverAddr, ok := os.LookupEnv("serverAddr")
	if !ok{
		log.Fatalln("serverAddr is not written in credentials.sh file")
	}

	conn, err := grpc.Dial(serverAddr, opts...)
	if err != nil {
		log.Fatalln(err.Error())
	}
	go func() {
		signs := make(chan os.Signal, 1)
		signal.Notify(signs, syscall.SIGINT, syscall.SIGTERM)
		sign := <-signs
		fmt.Println(sign)
		conn.Close()
		os.Exit(0)
	}()
	return conn
}

func NewPhotoClient() proto.NewPhotosClient {
	return proto.NewNewPhotosClient(conn)
}
func NewAuthClient() proto.AuthenticationClient {
	return proto.NewAuthenticationClient(conn)
}
