package proxy

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/zrpc"
)

type Interceptor struct {
}

func (i Interceptor) Do(ctx context.Context, conn *zrpc.RpcClient, fullMethodName string) error {
	fmt.Println("---> ", fullMethodName)
	return nil
}
