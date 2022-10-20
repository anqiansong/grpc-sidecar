package main

import (
	"context"
	"fmt"
	"log"

	"github.com/anqiansong/grpc-sidecar/example/grpc-proxy/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func main() {
	conn, err := grpc.Dial("127.0.0.1:9002",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
			md, _ := metadata.FromOutgoingContext(ctx)
			if len(md) != 0 {
				values := md.Get("grpc-proxy-date")
				if len(values) > 0 {
					fmt.Println("grpc-proxy-date:", values[0])
				}
			}

			return invoker(ctx, method, req, reply, cc, opts...)
		}))
	if err != nil {
		log.Fatalln(err)
	}

	client := pb.NewGreetServiceClient(conn)
	ctx := context.Background()
	resp, err := client.Echo(ctx, &pb.Req{
		In: "hello",
	})
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(resp.Out)
}
