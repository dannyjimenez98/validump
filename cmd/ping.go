package cmd

import (
	"context"
	"fmt"

	"time"

	"github.com/dannyjimenez98/validump/internal/databases/postgres"
	"github.com/spf13/cobra"
)

// newPingCmd represents the ping command
func newPingCmd(cfg *App) *cobra.Command {
	var pingCmd = &cobra.Command{
		Use:   "ping <db-alias>",
		Short: "test the database connection",
		Long:  ``,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dbalias := args[0]

			connection, err := postgres.Connect(cmd.Context(), cfg.ViperConfig, dbalias)
			if err != nil {
				return fmt.Errorf("create database connection: %w", err)
			}
			defer connection.Close(cmd.Context())

			ctx, cancel := context.WithTimeout(cmd.Context(), 5*time.Second)
			defer cancel()

			if err := connection.Ping(ctx); err != nil {
				return fmt.Errorf("ping database: %w", err)
			}

			fmt.Printf("DB '%s' successfully pinged!\n", dbalias)

			return nil
		},
	}
	return pingCmd
}
