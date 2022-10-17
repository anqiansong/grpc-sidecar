package main

import (
	"context"
	"log"

	"github.com/anqiansong/grpc-sidecar/example/pb"
	"github.com/anqiansong/grpc-sidecar/proxy"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

func main() {
	var address = "127.0.0.1:8080"
	var serviceName = "example.rpc"
	if err := proxy.RegisterSrv(proxy.Server{
		Name:    serviceName,
		Address: address,
		Timeout: 5000,
	}); err != nil {
		log.Fatalln(err)
	}

	server := zrpc.MustNewServer(zrpc.RpcServerConf{
		ServiceConf: service.ServiceConf{
			Name: serviceName,
		},
		ListenOn: address,
	}, func(server *grpc.Server) {
		pb.RegisterExampleServiceServer(server, &exampleServer{})
	})
	defer server.Stop()
	server.Start()

}

type exampleServer struct {
	pb.UnimplementedExampleServiceServer
}

func (e *exampleServer) Echo(_ context.Context, req *pb.ExampleReq) (*pb.ExampleResp, error) {
	return &pb.ExampleResp{
		Out: req.In,
	}, nil
}
