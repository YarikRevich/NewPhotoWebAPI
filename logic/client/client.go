package client

import (
	"NewPhotoWeb/log"
	"NewPhotoWeb/logic/proto"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
)

var (
	NewPhotoClient, NewPhotoAuthClient = New()
)

type Client struct {
	conn *grpc.ClientConn
}

func (c *Client) NewConnection() {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(50*10e6),
			grpc.MaxCallSendMsgSize(50*10e6),
		),
	}

	serverAddr, ok := os.LookupEnv("serverAddr")
	if !ok {
		log.Logger.Fatalln("serverAddr is not written in credentials.sh file")
	}

	conn, err := grpc.Dial(serverAddr, opts...)
	if err != nil {
		log.Logger.Fatalln(err.Error())
	}
	go func() {
		signs := make(chan os.Signal, 1)
		signal.Notify(signs, syscall.SIGINT, syscall.SIGTERM)
		sign := <-signs
		fmt.Println(sign)
		if err := conn.Close(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		os.Exit(0)
	}()
	c.conn = conn
}

func (c Client) GetNewPhotoClient() proto.NewPhotosClient {
	return proto.NewNewPhotosClient(c.conn)
}

func (c Client) GetNewPhotoAuthClient() proto.AuthenticationClient {
	return proto.NewAuthenticationClient(c.conn)
}

func New() (proto.NewPhotosClient, proto.AuthenticationClient) {
	l := new(Client)
	l.NewConnection()
	return l.GetNewPhotoClient(), l.GetNewPhotoAuthClient()
}

func Restart() {
	for {
		NewPhotoClient, NewPhotoAuthClient = New()
		if r, err := NewPhotoClient.Ping(context.Background(), &proto.PingRequest{}); err == nil && r.GetPong() {
			break
		}
		time.Sleep(5 * time.Second)
		log.Logger.ClientError()
	}
}
