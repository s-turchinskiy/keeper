package cmds

import (
	"github.com/s-turchinskiy/keeper/internal/client/models"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Run:   createVersionHandler(),
}

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register new user",
	Run:   withErrorHandling(createRegisterHandler()),
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login existing user",
	Run:   withErrorHandling(createLoginHandler()),
}

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync with remote repository",
	Run:   withErrorHandling(createSyncHandler()),
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new secret",
}

var addPasswordCmd = createSecretAddCommand(models.SecretTypePassword)
var addTextCmd = createSecretAddCommand(models.SecretTypeText)
var addBinaryCmd = createSecretAddCommand(models.SecretTypeBinary)
var addCardCmd = createSecretAddCommand(models.SecretTypeCard)

var getCmd = &cobra.Command{
	Use:   "get [uuid]",
	Short: "Get secret by UUID",
	Args:  cobra.ExactArgs(1),
	Run:   withErrorHandling(createSecretGetCommand()),
}

var deleteCmd = &cobra.Command{
	Use:   "delete [uuid]",
	Short: "Delete secret by UUID",
	Args:  cobra.ExactArgs(1),
	Run:   withErrorHandling(createSecretDeleteCommand()),
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all secrets",
	Run:   withErrorHandling(createSecretsListCommand()),
}
