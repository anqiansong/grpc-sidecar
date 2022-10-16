package proxy

import (
	"github.com/mwitkow/grpc-proxy/proxy"
	"google.golang.org/grpc"
	"net"
)

func NewServer(conf Config) error {
	th := proxy.TransparentHandler(director)
	option := grpc.UnknownServiceHandler(th)
	s := grpc.NewServer(option)
	l, err := net.Listen("tcp", conf.GrpcServer.Address)
	if err != nil {
		return err
	}

	defer s.GracefulStop()
	return s.Serve(l)
}
