package proxy

import (
	"context"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Interceptor struct {
}

func NewInterceptor() *Interceptor {
	return &Interceptor{}
}

func (i Interceptor) Do(ctx context.Context, conn *zrpc.RpcClient, fullMethodName string, config Config) error {
	for _, v := range config.Interceptor {
		if fullMethodName == v.Method && v.Enable {
			return status.Errorf(codes.Aborted, "aborted")
		}
	}
	return nil
}
