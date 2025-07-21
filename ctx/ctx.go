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

func CloneContextWithValues(ctx context.Context, keys []interface{}) context.Context {
	newCtx := context.Background()
	for _, key := range keys {
		val := ctx.Value(key)
		if val != nil {
			newCtx = context.WithValue(newCtx, key, val)
		}
	}
	return newCtx
}
