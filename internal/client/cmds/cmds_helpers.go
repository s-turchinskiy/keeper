package cmds

import (
	"fmt"
	"github.com/s-turchinskiy/keeper/internal/client/service"
	"log"

	"github.com/spf13/cobra"
)

func setFlags() {
	addPasswordCmd.Flags().String("name", "", "Secret name (required)")
	addPasswordCmd.Flags().String("username", "", "Username (required)")
	addPasswordCmd.Flags().String("password", "", "Password (required)")
	addPasswordCmd.Flags().String("url", "", "URL (optional)")
	addPasswordCmd.Flags().String("metadata", "", "Metadata (optional)")
	markFlagsRequired(addPasswordCmd, "name", "username", "password")

	addTextCmd.Flags().String("name", "", "Secret name (required)")
	addTextCmd.Flags().String("content", "", "Text content (required)")
	addTextCmd.Flags().String("metadata", "", "Metadata (optional)")
	markFlagsRequired(addTextCmd, "name", "content")

	addBinaryCmd.Flags().String("name", "", "Secret name (required)")
	addBinaryCmd.Flags().String("file", "", "File path (required)")
	addBinaryCmd.Flags().String("metadata", "", "Metadata (optional)")
	markFlagsRequired(addBinaryCmd, "name", "file")

	addCardCmd.Flags().String("name", "", "Secret name (required)")
	addCardCmd.Flags().String("number", "", "Card number (required)")
	addCardCmd.Flags().String("holder", "", "Card holder name (required)")
	addCardCmd.Flags().String("expiry", "", "Expiry date (required)")
	addCardCmd.Flags().String("cvv", "", "CVV code (required)")
	addCardCmd.Flags().String("metadata", "", "Metadata (optional)")
	markFlagsRequired(addCardCmd, "name", "number", "holder", "expiry", "cvv")

	getCmd.Flags().Bool("full", false, "Show all data including passwords/CVV")
	getCmd.Flags().String("export", "", "Export to file path")

	addCmd.AddCommand(addPasswordCmd)
	addCmd.AddCommand(addTextCmd)
	addCmd.AddCommand(addBinaryCmd)
	addCmd.AddCommand(addCardCmd)
}

func markFlagsRequired(cmd *cobra.Command, flags ...string) {
	for _, flag := range flags {
		if err := cmd.MarkFlagRequired(flag); err != nil {
			log.Fatalf("Error marking flag %s required: %v", flag, err)
		}
	}
}

func addCommands(rootCmd *cobra.Command) {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(registerCmd)
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(syncCmd)
}

func getServiceFromCommand(cmd *cobra.Command) service.Servicer {
	rootCmd := cmd.Root()
	srvc, ok := rootCmd.Context().Value(serviceContextKey).(service.Servicer)
	if !ok {
		log.Fatal("app not found in command context")
	}
	return srvc
}

func getStringFlag(cmd *cobra.Command, name string) string {
	value, _ := cmd.Flags().GetString(name)
	return value
}

func withErrorHandling(fn func(cmd *cobra.Command, args []string) error) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		err := fn(cmd, args)
		short := cmd.Short
		if err != nil {
			fmt.Printf("Failed to %s: %v\n", short, err)
			return
		}
	}
}
