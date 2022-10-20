package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/anqiansong/grpc-sidecar/example/grpc-proxy/pb"
	"github.com/mwitkow/grpc-proxy/proxy"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func main() {
	go listenProxy()
	listenServer()
}

func listenProxy() {
	th := proxy.TransparentHandler(func(ctx context.Context, fullMethodName string) (context.Context, *grpc.ClientConn, error) {
		fmt.Println("-----")
		conn, err := grpc.Dial("127.0.0.1:9001", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return nil, nil, err
		}

		newCtx := metadata.NewOutgoingContext(ctx, metadata.New(map[string]string{
			"grpc-proxy-date": time.Now().Format("2006-01-02 15:04:05"),
		}))
		return newCtx, conn, nil
	})

	serviceHandlerOption := grpc.UnknownServiceHandler(th)
	l, err := net.Listen("tcp", ":9002")
	if err != nil {
		log.Fatalln(err)
	}

	server := grpc.NewServer(serviceHandlerOption)
	defer server.GracefulStop()
	server.Serve(l)
}

func listenServer() {
	l, err := net.Listen("tcp", ":9001")
	if err != nil {
		log.Fatalln(err)
	}

	server := grpc.NewServer()
	pb.RegisterGreetServiceServer(server, &Greet{})
	defer server.GracefulStop()
	server.Serve(l)
}

type Greet struct {
	pb.UnimplementedGreetServiceServer
}

func (g *Greet) Echo(ctx context.Context, req *pb.Req) (*pb.Resp, error) {
	return &pb.Resp{Out: req.In}, nil
}
