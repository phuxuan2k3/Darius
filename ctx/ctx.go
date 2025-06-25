package ctx

import (
	"context"

	"google.golang.org/grpc/metadata"
)

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
