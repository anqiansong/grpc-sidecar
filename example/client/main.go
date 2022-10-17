package main

import (
	"context"
	"fmt"
	"time"

	"github.com/anqiansong/grpc-sidecar/example/pb"
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/core/logx"
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
	for {
		time.Sleep(5 * time.Second)
		resp, err := c.Echo(ctx, &pb.ExampleReq{
			In: fmt.Sprintf("hello from %v", time.Now().Format("2006-01-02 15:04:05")),
		})
		if err != nil {
			logx.Error(err)
			continue
		}
		fmt.Println(resp.Out)
	}
}
