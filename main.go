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
	if err := proxy.ListenSrv(ch,
		proxy.NewLimiter(),
		proxy.NewInterceptor(),
		proxy.NewAuth(),
	); err != nil {
		log.Fatalln(err)
	}
}
