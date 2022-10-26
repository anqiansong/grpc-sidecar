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
	conn, err := grpc.Dial("10.211.55.4:8000",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln(err)
	}

	client := pb.NewGreetServiceClient(conn)
	ctx := context.Background()
	var header, trailer metadata.MD
	resp, err := client.Echo(ctx, &pb.Req{
		In: "hello",
	}, grpc.Header(&header), grpc.Trailer(&trailer))
	if err != nil {
		log.Fatalln(err)
	}

	if header != nil {
		for k, v := range header {
			fmt.Println(k, v)
		}
	}
	if trailer != nil {
		for k, v := range trailer {
			fmt.Println(k, v)
		}
	}
	fmt.Println(resp.Out)
}
