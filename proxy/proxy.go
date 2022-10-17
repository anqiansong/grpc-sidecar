package proxy

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"sync"

	"github.com/mwitkow/grpc-proxy/proxy"
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type Proxy struct {
	clientManager sync.Pool
	filters       []Filter
	once          sync.Once
	ch            *ConfigHandler
	conf          *Config
	lock          sync.Mutex
}

func NewProxy(ch *ConfigHandler, filters ...Filter) *Proxy {
	var p = &Proxy{
		filters: filters,
		ch:      ch,
	}
	return p
}

func (p *Proxy) SyncConfig() {
	p.ch.Listen(func(conf *Config) {
		p.lock.Lock()
		p.conf = conf
		p.lock.Unlock()
	})
}

func (p *Proxy) Run(server Server) {
	logx.Info("enter proxy")
	p.once.Do(func() {
		logx.Info("run proxy")
		p.clientManager = sync.Pool{
			New: func() any {
				client, err := zrpc.NewClient(zrpc.RpcClientConf{
					Target: server.Address,
				})
				if err != nil {
					return nil
				}
				return client
			},
		}

		host, port, err := net.SplitHostPort(server.Address)
		if err != nil {
			log.Fatalln(err)
		}

		portInt, err := strconv.ParseUint(port, 10, 64)
		if err != nil {
			log.Fatalf("parse port %q error: %+v", port, err)
		}

		proxyServerAddress := net.JoinHostPort(host, fmt.Sprint(portInt+10))
		server := zrpc.MustNewServer(zrpc.RpcServerConf{
			ServiceConf: service.ServiceConf{
				Name: server.Name,
			},
			ListenOn: proxyServerAddress,
			Etcd: discov.EtcdConf{
				Hosts: []string{"127.0.0.1:2379"},
				Key:   server.Name,
			},
			Timeout: server.Timeout,
		}, func(server *grpc.Server) {})
		th := proxy.TransparentHandler(p.director)
		serviceHandlerOption := grpc.UnknownServiceHandler(th)
		server.AddOptions(serviceHandlerOption)
		defer server.Stop()
		server.Start()
	})
}

func (p *Proxy) director(ctx context.Context, fullMethodName string) (context.Context, *grpc.ClientConn, error) {
	v := p.clientManager.Get()
	err, ok := v.(error)
	if ok {
		return ctx, nil, err
	}

	conn, ok := v.(*zrpc.RpcClient)
	if !ok {
		return ctx, nil, errors.New("get client conn instance failed")
	}

	for _, f := range p.filters {
		err = f.Do(ctx, conn, fullMethodName, p.conf)
		if err != nil {
			return ctx, conn.Conn(), err
		}
	}

	return ctx, conn.Conn(), nil
}
