package models

import (
	"encoding/json"
	"github.com/s-turchinskiy/keeper/models/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/s-turchinskiy/keeper/internal/client/crypto"
)

func ConvertLocalSecretToRemoteSecret(cryptor crypto.Cryptor, localSecret *LocalSecret) (*RemoteSecret, error) {
	secretData, err := localSecret.ParseData()
	if err != nil {
		return nil, err
	}

	secretDataContainer := &SecretDataContainer{
		Type:       localSecret.Type,
		Name:       localSecret.Name,
		SecretData: secretData,
	}

	remoteData, err := json.Marshal(secretDataContainer)
	if err != nil {
		return nil, err
	}

	encryptedRemoteData, err := cryptor.EncryptSecretData(remoteData)
	if err != nil {
		return nil, err
	}

	remoteSecret := &RemoteSecret{
		Name:         localSecret.Name,
		LastModified: localSecret.LastModified,
		Hash:         localSecret.Hash,
		Data:         encryptedRemoteData,
	}

	return remoteSecret, nil
}

func ConvertRemoteSecretToLocalSecret(cryptor crypto.Cryptor, remoteSecret *RemoteSecret) (*LocalSecret, error) {
	remoteDecryptedData, err := cryptor.DecryptSecretData(remoteSecret.Data)
	if err != nil {
		return nil, err
	}

	var secretDataContainer SecretDataContainer
	if err := json.Unmarshal(remoteDecryptedData, &secretDataContainer); err != nil {
		return nil, err
	}

	localSecret := &LocalSecret{
		Name:         remoteSecret.Name,
		Type:         secretDataContainer.Type,
		LastModified: remoteSecret.LastModified,
		Hash:         remoteSecret.Hash,
	}

	err = localSecret.SetData(cryptor, secretDataContainer.SecretData)
	if err != nil {
		return nil, err
	}

	return localSecret, nil
}

func ConvertProtoSecretToRemoteSecret(secretResp *proto.Secret) *RemoteSecret {
	return &RemoteSecret{
		Name:         secretResp.GetId(),
		LastModified: secretResp.GetLastModified().AsTime(),
		Hash:         secretResp.GetHash(),
		Data:         secretResp.GetData(),
	}
}

func ConvertRemoteSecretToProtoSecret(secret *RemoteSecret) *proto.Secret {
	return &proto.Secret{
		Id:           secret.Name,
		LastModified: timestamppb.New(secret.LastModified),
		Hash:         secret.Hash,
		Data:         secret.Data,
	}
}
