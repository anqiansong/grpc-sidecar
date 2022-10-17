package proxy

import (
	"context"

	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

func ListenControlPanel(ch *ConfigHandler) error {
	var cs = NewCPServer(ch)
	ch.wg.Add(1)
	server, err := zrpc.NewServer(zrpc.RpcServerConf{
		ServiceConf: service.ServiceConf{
			Name: "cp-rpc",
		},
		ListenOn: cpListen,
		Timeout:  5000,
	}, func(server *grpc.Server) {
		RegisterCPServiceServer(server, cs)
	})
	if err != nil {
		return err
	}

	defer server.Stop()
	server.Start()
	return nil
}

type cpServer struct {
	UnimplementedCPServiceServer
	ch *ConfigHandler
}

func NewCPServer(ch *ConfigHandler) *cpServer {
	return &cpServer{
		ch: ch,
	}
}

func (c cpServer) SyncConfig(_ context.Context, request *CPRequest) (*CPResponse, error) {
	var conf = configTransfer(request)
	c.ch.Put(conf)
	return &CPResponse{}, nil
}

func configTransfer(request *CPRequest) Config {
	return Config{}
}
