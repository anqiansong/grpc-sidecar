package proxy

import (
	"context"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type Auth struct {
}

type Route struct {
	Method string `json:"method"`
	Enable bool   `json:"enable"`
}

type AuthConfig struct {
	Enable bool    `json:"enable"` // 鉴定
	Route  []Route `json:"route"`
}

func (ac AuthConfig) SelectRoute(fullMethodName string) Route {
	for _, v := range ac.Route {
		if v.Method == fullMethodName {
			return v
		}
	}

	return Route{
		Method: fullMethodName,
	}
}

func NewAuth() *Auth {
	return &Auth{}
}

func (a Auth) Do(ctx context.Context, conn *zrpc.RpcClient, fullMethodName string, config Config) error {
	if config.AuthConfig.Enable {
		if err := a.Authentication(ctx, config.AuthConfig.SelectRoute(fullMethodName)); err != nil {
			return err
		}
	}

	return nil
}

func (a Auth) Authentication(ctx context.Context, route Route) error {
	if !route.Enable {
		return nil
	}

	// sample
	values, ok := metadata.FromIncomingContext(ctx)
	if !ok || len(values) == 0 {
		return status.Error(codes.Unauthenticated, "missing token")
	}

	token := values.Get("appID")
	if len(token) == 0 {
		return status.Error(codes.Unauthenticated, "missing token")
	}

	// mock
	if token[0] != "123456" {
		return status.Error(codes.Unauthenticated, "auth failed")
	}

	return nil
}
