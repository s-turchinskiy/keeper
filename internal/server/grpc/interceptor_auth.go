package grpc

import (
	"context"
	"github.com/s-turchinskiy/keeper/internal/server/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type contextKey string

const userIDKey contextKey = "userID"

func AuthInterceptor(service *service.Service) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if info.FullMethod == "/keeper.AuthService/Login" ||
			info.FullMethod == "/keeper.AuthService/Register" ||
			info.FullMethod == "/keeper.AuthService/GetConnectionNumber" {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}
		tokens := md["authorization"]
		if len(tokens) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing authorization token")
		}

		token := tokens[0]
		userID, err := service.TokenManager.ValidateToken(token)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}

		ctx = context.WithValue(ctx, userIDKey, userID)
		return handler(ctx, req)
	}
}

func getUserIDFromContext(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(userIDKey).(string)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "user not found in context")
	}

	return userID, nil
}
