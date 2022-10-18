package main

import (
	"fmt"
	"log"

	"github.com/anqiansong/grpc-sidecar/proxy"
)

func main() {
	var ch = proxy.NewConfigHandler()
	go func() {
		if err := proxy.ListenControlPanel(ch); err != nil {
			log.Fatalln(err)
		}
	}()

	defer ch.Close()
	ch.Wait()

	fmt.Println("started sidecar ...")

	p := proxy.NewProxy(ch,
		proxy.NewLimiter(),
		proxy.NewInterceptor(),
		proxy.NewAuth())
	p.RunClientProxy()

	if err := proxy.ListenSrv(p); err != nil {
		log.Fatalln(err)
	}
}
