package main

import (
	"context"
	"fmt"
	"time"

	"github.com/anqiansong/grpc-sidecar/example/pb"
	"github.com/anqiansong/grpc-sidecar/proxy"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc/metadata"
)

func main() {
	client := zrpc.MustNewClient(zrpc.RpcClientConf{
		Target: "127.0.0.1:9000",
	})
	c := pb.NewExampleServiceClient(client.Conn())
	ctx := context.Background()
	for {
		// 10 qps
		time.Sleep(100 * time.Millisecond)
		newCtx := metadata.NewOutgoingContext(ctx, metadata.New(map[string]string{
			proxy.GrpcServerName: "example.rpc",
			"appID":              "123456",
		}))
		resp, err := c.Echo(newCtx, &pb.ExampleReq{
			In: fmt.Sprintf("hello from %v", time.Now().Format("2006-01-02 15:04:05")),
		})
		if err != nil {
			logx.Error(err)
			continue
		}
		fmt.Println(resp.Out)
	}
}
