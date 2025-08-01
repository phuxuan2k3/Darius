package ctx

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const HttpCodeHeader string = "X-Http-Code"

func GetUserIdFromContext(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", nil
	}

	userID := md.Get("x-user-id")

	if len(userID) == 0 {
		return "", nil
	}
	return userID[0], nil
}

func SetHeaders(ctx context.Context, key, value string) error {
	md := metadata.Pairs(key, value)
	return grpc.SetHeader(ctx, md)
}

func CloneContextWithValues(ctx context.Context) context.Context {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.MD{}
	}
	clonedCtx := metadata.NewIncomingContext(context.Background(), md)
	return clonedCtx
}
