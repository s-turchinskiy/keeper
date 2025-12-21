package grpc

import (
	"github.com/s-turchinskiy/keeper/internal/server/models"
	"github.com/s-turchinskiy/keeper/models/proto"
	"github.com/s-turchinskiy/keeper/pkd/ternary"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertProtoSecretToServerSecret(grpcSecret *proto.Secret, userID string) *models.Secret {
	return &models.Secret{
		ID:           grpcSecret.GetId(),
		UserID:       userID,
		Data:         grpcSecret.GetData(),
		Hash:         grpcSecret.GetHash(),
		LastModified: grpcSecret.GetLastModified().AsTime(),
		Deleted:      ternary.Bool(grpcSecret.Deleted == nil, false, grpcSecret.Deleted),
	}
}

func convertServerSecretToProtoSecret(secret *models.Secret) *proto.Secret {
	return &proto.Secret{
		Id:           secret.ID,
		Data:         secret.Data,
		Hash:         secret.Hash,
		LastModified: timestamppb.New(secret.LastModified),
		Deleted:      &secret.Deleted,
	}
}
