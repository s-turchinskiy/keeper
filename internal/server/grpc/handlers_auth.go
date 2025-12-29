package grpc

import (
	"context"
	"errors"
	"github.com/s-turchinskiy/keeper/internal/server/service"
	"github.com/s-turchinskiy/keeper/models/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthHandler struct {
	proto.UnimplementedAuthServiceServer
	service service.Servicer
}

func (h *AuthHandler) GetConnectionNumber(ctx context.Context, req *proto.GetConnectionNumberRequest) (*proto.GetConnectionNumberResponse, error) {
	connectionNumber := h.service.GetNewConnectionNumber(ctx)

	resp := &proto.GetConnectionNumberResponse{
		ConnectionNumber: connectionNumber,
	}

	return resp, nil
}

func NewAuthHandler(service service.Servicer) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) Register(ctx context.Context, req *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	user, err := h.service.Register(ctx, req.GetLogin(), req.GetPassword())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	resp := &proto.RegisterResponse{
		UserId: user.ID,
	}

	return resp, nil
}

func (h *AuthHandler) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	token, user, err := h.service.Login(ctx, req.GetLogin(), req.GetPassword())
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserNotFound):
			return nil, status.Error(codes.NotFound, "user not found")
		case errors.Is(err, service.ErrInvalidCredentials):
			return nil, status.Error(codes.Unauthenticated, "invalid password")
		default:
			return nil, status.Error(codes.Internal, "internal error")
		}
	}
	resp := &proto.LoginResponse{
		Token:  token,
		UserId: user.ID,
	}

	return resp, nil
}
