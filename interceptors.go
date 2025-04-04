package krkstops

import (
	"context"
	"errors"
	"math/rand/v2"

	"google.golang.org/grpc"
)

func InjectFailure(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	if rand.IntN(2) < 1 {
		errorStr := []string{
			"stop not found",
			"got non 200 status code",
			"request failed",
		}
		return nil, errors.New(errorStr[rand.IntN(len(errorStr))])
	} else {
		return handler(ctx, req)
	}
}
