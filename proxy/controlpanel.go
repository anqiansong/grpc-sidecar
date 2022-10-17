package proxy

import (
	"context"
	"encoding/json"

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
	conf, err := configTransfer(request)
	if err != nil {
		return nil, err
	}

	c.ch.Put(conf)
	return &CPResponse{}, nil
}

func configTransfer(request *CPRequest) (*Config, error) {
	var conf Config
	if err := json.Unmarshal(request.In, &conf); err != nil {
		return nil, err
	}
	return &conf, nil
}
