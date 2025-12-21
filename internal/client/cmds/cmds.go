package cmds

import (
	"context"
	"github.com/s-turchinskiy/keeper/internal/client/service"
	"github.com/spf13/cobra"
)

type contextKey string

const (
	serviceContextKey contextKey = "app"
)

type CobraCommand struct {
	rootCmd *cobra.Command
}

func New(service service.Servicer) *CobraCommand {
	rootCmd := &cobra.Command{
		Use:   "keeper",
		Short: "Zero-Knowledge secret manager",
	}

	ctx := context.WithValue(context.Background(), serviceContextKey, service)
	rootCmd.SetContext(ctx)

	setFlags()
	addCommands(rootCmd)
	return &CobraCommand{
		rootCmd: rootCmd,
	}
}

func (c *CobraCommand) Run() error {
	return c.rootCmd.Execute()
}
