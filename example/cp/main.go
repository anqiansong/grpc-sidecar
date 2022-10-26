package main

import (
	"context"
	"flag"
	"io/ioutil"
	"log"

	"github.com/anqiansong/grpc-sidecar/proxy"
	"github.com/zeromicro/go-zero/zrpc"
)

var f = flag.String("f", "config.json", "config file")

func main() {
	flag.Parse()

	data, err := ioutil.ReadFile(*f)
	if err != nil {
		log.Fatalln(err)
	}

	c := zrpc.MustNewClient(zrpc.RpcClientConf{
		Target: "127.0.0.1:9010",
	})
	client := proxy.NewCPServiceClient(c.Conn())
	_, err = client.SyncConfig(context.Background(), &proxy.CPRequest{
		In: data,
	})
	if err != nil {
		log.Fatalln(err)
	}
}
