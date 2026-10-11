package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/cobra"
)

// pingCmd represents the ping command
var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "test the database connection",
	Long: ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		connStr := os.Getenv("DB_CONNECTION_STR")
		if connStr == "" {
			return errors.New("database connection string not set")
		}

		ctx, cancel := context.WithTimeout(cmd.Context(), 5*time.Second)
		defer cancel()

		connection, err := pgxpool.New(ctx, connStr)
		if err != nil {
			return fmt.Errorf("create database connection: %w", err) 
		}
		defer connection.Close()

		if err := connection.Ping(ctx); err != nil {
			return fmt.Errorf("ping database: %w", err)
		}

		fmt.Println("DB connection successful!")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(pingCmd)
}
