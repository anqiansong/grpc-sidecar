package main

import (
	"context"
	"fmt"
	"log"

	"github.com/anqiansong/grpc-sidecar/example/pb"
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/zrpc"
)

func main() {
	client := zrpc.MustNewClient(zrpc.RpcClientConf{
		Etcd: discov.EtcdConf{
			Hosts: []string{"127.0.0.1:2379"},
			Key:   "example.rpc",
		},
	})
	c := pb.NewExampleServiceClient(client.Conn())
	ctx := context.Background()
	resp, err := c.Echo(ctx, &pb.ExampleReq{
		In: "hello world",
	})
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(resp.Out)
}
