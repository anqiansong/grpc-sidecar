package proxy

import (
	"context"

	"github.com/zeromicro/go-zero/core/limit"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Limiter struct {
	r *redis.Redis
}

func NewLimiter() *Limiter {
	r := redis.New("127.0.0.1:6379")
	return &Limiter{
		r: r,
	}
}

func (l Limiter) Do(ctx context.Context, conn *zrpc.RpcClient, fullMethodName string, config Config) error {
	if !config.Limiter.Enable {
		return nil
	}

	quota := config.Limiter.Qps
	if quota <= 0 {
		quota = 5
	}

	limiter := limit.NewPeriodLimit(1, quota, l.r, "")
	val, err := limiter.Take(fullMethodName)
	if err != nil {
		return err
	}
	switch val {
	case limit.Allowed:
		return nil
	default:
		return status.Error(codes.Unavailable, "qps overflow")
	}
}
