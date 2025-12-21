package cmds

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/s-turchinskiy/keeper/internal/client/models"
	"github.com/s-turchinskiy/keeper/internal/utils/buildinfo"
	"github.com/s-turchinskiy/keeper/internal/utils/filecheckerutils"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

func createVersionHandler() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		fmt.Printf("Build version: %s\n", buildinfo.BuildVersion)
		fmt.Printf("Build date: %s\n", buildinfo.BuildDate)
		fmt.Printf("Build commit: %s\n", buildinfo.BuildCommit)
	}
}

func createRegisterHandler() func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		service := getServiceFromCommand(cmd)

		login := getStringFlag(cmd, "login")
		password := getStringFlag(cmd, "password")

		err := service.Register(context.Background(), login, password)
		if err != nil {
			return err
		}
		return nil
	}
}

func createLoginHandler() func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		service := getServiceFromCommand(cmd)

		login := getStringFlag(cmd, "login")
		password := getStringFlag(cmd, "password")

		err := service.Login(context.Background(), login, password)
		if err != nil {
			return err
		}
		return nil
	}
}

func createSyncHandler() func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		service := getServiceFromCommand(cmd)
		err := service.SyncSecrets(context.Background())
		if err != nil {
			return err
		}
		return nil
	}
}

func createSecretAddCommand(secretType string) *cobra.Command {
	return &cobra.Command{
		Use:   secretType,
		Short: fmt.Sprintf("Add %s secret", secretType),
		Run:   withErrorHandling(createSecretHandler(secretType)),
	}
}

func createSecretHandler(secretType string) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		base := models.BaseSecret{
			Type:     secretType,
			Name:     getStringFlag(cmd, "name"),
			Metadata: getStringFlag(cmd, "metadata"),
		}

		var data models.SecretData
		switch secretType {
		case models.SecretTypePassword:
			data = models.LoginData{
				Username: getStringFlag(cmd, "username"),
				Password: getStringFlag(cmd, "password"),
				URL:      getStringFlag(cmd, "url"),
			}
		case models.SecretTypeText:
			data = models.TextData{
				Content: getStringFlag(cmd, "content"),
			}
		case models.SecretTypeBinary:
			filePath := getStringFlag(cmd, "file")
			checker := filecheckerutils.NewFileChecker()
			if err := checker.CheckFileSize(filePath, models.MaxFileSize); err != nil {
				return err
			}
			content, err := os.ReadFile(filePath)
			if err != nil {
				return fmt.Errorf("failed to read file: %w", err)
			}
			data = models.FileData{
				FileName: filepath.Base(filePath),
				FileSize: int64(len(content)),
				Content:  base64.StdEncoding.EncodeToString(content),
			}
		case models.SecretTypeCard:
			data = models.CardData{
				Number: getStringFlag(cmd, "number"),
				Holder: getStringFlag(cmd, "holder"),
				Expiry: getStringFlag(cmd, "expiry"),
				CVV:    getStringFlag(cmd, "cvv"),
			}
		default:
			return fmt.Errorf("unsupported secret type: %s", secretType)
		}

		service := getServiceFromCommand(cmd)

		createdSecret, err := service.CreateSecret(context.Background(), base, data)
		if err != nil {
			return err
		}

		err = displaySecret(createdSecret, false)
		if err != nil {
			return err
		}

		return nil
	}
}

func createSecretGetCommand() func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		uuid := args[0]
		full, _ := cmd.Flags().GetBool("full")

		service := getServiceFromCommand(cmd)

		secret, err := service.ReadSecret(context.Background(), uuid)
		if err != nil {
			return err
		}

		err = displaySecret(secret, full)
		if err != nil {
			return err
		}

		return nil
	}
}

func createSecretDeleteCommand() func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		uuid := args[0]

		service := getServiceFromCommand(cmd)
		_, err := service.DeleteSecret(context.Background(), uuid)
		if err != nil {
			return err
		}

		return nil
	}
}

func createSecretsListCommand() func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {

		service := getServiceFromCommand(cmd)
		secrets, err := service.ListLocalSecrets(context.Background())
		if err != nil {
			return err
		}

		displaySecrets(secrets)

		return nil
	}
}
