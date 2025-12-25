package grpc

import (
	"context"
	"errors"
	"github.com/s-turchinskiy/keeper/internal/server/repository/postgres"
	"github.com/s-turchinskiy/keeper/internal/server/service"
	"github.com/s-turchinskiy/keeper/models/proto"
	"google.golang.org/grpc"
	"log"

	"github.com/s-turchinskiy/keeper/internal/server/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type SecretHandler struct {
	proto.UnimplementedSecretServiceServer
	service             service.Servicer
	secretsForClientsCh chan []*models.Secret
}

func NewSecretHandler(service service.Servicer) *SecretHandler {
	return &SecretHandler{
		service:             service,
		secretsForClientsCh: make(chan []*models.Secret),
	}
}

func (h *SecretHandler) SetSecret(ctx context.Context, req *proto.SetSecretRequest) (*proto.SetSecretResponse, error) {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "authentication required")
	}

	reqSecret := req.GetSecret()

	secret := convertProtoSecretToServerSecret(reqSecret, userID)

	err = h.service.SetSecret(ctx, secret)
	if err != nil {
		log.Printf("SetSecret failed: %v", err)
		return nil, status.Error(codes.Internal, "failed to set secret, err: "+err.Error())
	}

	resp := &proto.SetSecretResponse{}
	resp.Success = true

	h.secretsForClientsCh <- []*models.Secret{secret}

	return resp, nil
}

func (h *SecretHandler) GetSecret(ctx context.Context, req *proto.GetSecretRequest) (*proto.GetSecretResponse, error) {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "authentication required")
	}

	secret, err := h.service.GetSecret(ctx, userID, req.GetSecretId())
	if err != nil {
		log.Printf("GetSecret failed: %v", err)
		return nil, status.Error(codes.Internal, "failed to get secret")
	}

	respSecret := &proto.Secret{
		Id:           secret.ID,
		Hash:         secret.Hash,
		LastModified: timestamppb.New(secret.LastModified),
		Data:         secret.Data,
	}

	resp := &proto.GetSecretResponse{}
	resp.Secret = respSecret

	return resp, nil
}

func (h *SecretHandler) DeleteSecret(ctx context.Context, req *proto.DeleteSecretRequest) (*proto.DeleteSecretResponse, error) {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "authentication required")
	}

	err = h.service.DeleteSecret(ctx, userID, req.GetSecretId())
	if err != nil && !errors.Is(err, postgres.ErrSecretNotFound) {
		log.Printf("DeleteSecret failed: %v", err)
		return nil, status.Error(codes.Internal, "failed to delete secret")
	}

	resp := &proto.DeleteSecretResponse{}
	resp.Success = true

	secret, err := h.service.GetSecret(ctx, userID, req.GetSecretId())
	if err == nil {
		h.secretsForClientsCh <- []*models.Secret{secret}
	} else {
		log.Printf("error send deleted secret in clients, user id: %s, name: %s, error : %v", userID, req.GetSecretId(), err)
	}

	return resp, nil
}

func (h *SecretHandler) ListSecrets(ctx context.Context, req *proto.ListSecretsRequest) (*proto.ListSecretsResponse, error) {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "authentication required")
	}

	secrets, err := h.service.ListSecrets(ctx, userID)
	if err != nil {
		log.Printf("ListSecrets failed: %v", err)
		return nil, status.Error(codes.Internal, "failed to retrieve secrets, err: "+err.Error())
	}

	respSecrets := make([]*proto.Secret, len(secrets))
	for i, secret := range secrets {
		respSecrets[i] = convertServerSecretToProtoSecret(secret)
	}

	resp := &proto.ListSecretsResponse{}
	resp.Secrets = respSecrets

	return resp, nil
}

func (h *SecretHandler) SyncSecretsFromClient(ctx context.Context, req *proto.SyncSecretsFromClientRequest) (*proto.SyncSecretsFromClientResponse, error) {

	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "authentication required")
	}

	reqSecrets := req.GetSecrets()

	secrets := make([]*models.Secret, len(reqSecrets))
	for _, reqSecret := range reqSecrets {
		secrets = append(secrets, convertProtoSecretToServerSecret(reqSecret, userID))
	}

	updateInClients, err := h.service.SyncFromClient(ctx, userID, secrets)
	if err != nil {
		log.Printf("SyncSecretsFromClient failed: %v", err)
		return nil, status.Error(codes.Internal, "failed to sync secrets from client, err: "+err.Error())
	}

	resp := &proto.SyncSecretsFromClientResponse{}
	resp.Success = true

	h.secretsForClientsCh <- updateInClients
	return resp, nil
}

// GetUpdatedSecrets Отправка измененных секретов всем подключенным клиентам
func (h *SecretHandler) GetUpdatedSecrets(req *proto.GetUpdatedSecretsRequest, g grpc.ServerStreamingServer[proto.GetUpdatedSecretsResponse]) error {

	//надо бы прерывать, но контекста когда stream нет

	//TODO: надо сделать проверку на userID, отправлять только на подключенные под именно этим юзером клиенты

	//TODO: это все работает ненадежно, то клиент получает данные стрима, то нет и непонятно почему
	for secretsForClients := range h.secretsForClientsCh {

		respSecrets := make([]*proto.Secret, len(secretsForClients))
		for i, secret := range secretsForClients {
			respSecrets[i] = convertServerSecretToProtoSecret(secret)
		}

		resp := &proto.GetUpdatedSecretsResponse{}
		resp.Secrets = respSecrets

		if err := g.Send(resp); err != nil {
			log.Printf("GetUpdatedSecrets failed: %v", err)
		}
	}

	return nil

}
