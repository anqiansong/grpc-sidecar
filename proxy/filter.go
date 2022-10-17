package proxy

import (
	"context"

	"github.com/zeromicro/go-zero/zrpc"
)

type Filter interface {
	Do(ctx context.Context, conn *zrpc.RpcClient, fullMethodName string, config *Config) error
}
