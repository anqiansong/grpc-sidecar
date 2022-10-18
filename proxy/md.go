package proxy

import (
	"context"

	"google.golang.org/grpc/metadata"
)

func FromIncomingContext(ctx context.Context, key string) string {
	md, _ := metadata.FromIncomingContext(ctx)
	if len(md) == 0 {
		return ""
	}
	val := md.Get(key)
	if len(val) == 0 {
		return ""
	}

	return val[0]
}
