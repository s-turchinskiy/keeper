package grpc

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func LoggingInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {

		var err error
		connNumber := "N/A"
		if info.FullMethod != "/keeper.AuthService/GetConnectionNumber" {
			connNumber, err = getConnectionNumber(ctx)

			if err != nil {
				log.Printf("gRPC request failed: connection number: %s, method: %s, error: %v",
					connNumber, info.FullMethod, err)
				return nil, err
			}
		}

		start := time.Now()

		resp, err := handler(ctx, req)

		duration := time.Since(start)
		statusCode := status.Code(err)

		if err != nil {
			log.Printf("gRPC request failed: connection number: %s, method: %s, duration: %s, status: %s, error: %v",
				connNumber, info.FullMethod, duration, statusCode, err)
		} else {
			log.Printf("gRPC request completed: connection number: %s, method: %s,  duration: %s, status: %s",
				connNumber, info.FullMethod, duration, statusCode)
		}

		return resp, err
	}
}

func getConnectionNumber(ctx context.Context) (string, error) {

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.InvalidArgument, "missing metadata")
	}
	dates := md["connectionnumber"]
	if len(dates) == 0 {
		return "", status.Error(codes.InvalidArgument, "missing connection number")
	}

	return dates[0], nil
}
