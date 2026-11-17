package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dannyjimenez98/validump/internal/config"
	"github.com/dannyjimenez98/validump/internal/databases/postgres"
	"github.com/spf13/cobra"
)

// newBackupCmd represents the backup command
func newBackupCmd(cfg *App) *cobra.Command {
	var backupCmd = &cobra.Command{
		Use:   "backup <db-alias>",
		Short: "Backup a database",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			dbalias := args[0]

			backupDirPath := cfg.ViperConfig.GetString("backup_dir")
			if backupDirPath == "" {
				home, err := os.UserHomeDir()
				if err != nil {
					return fmt.Errorf("get user home directory: %w", err)
				}
				backupDirPath = filepath.Join(home, "validump", dbalias)
			}

			if err := os.MkdirAll(backupDirPath, 0700); err != nil {
				return fmt.Errorf("create backups directory for database '%s': %w", dbalias, err)
			}

			dbCfg, err := config.LoadDB(cfg.ViperConfig, dbalias)
			if err != nil {
				return fmt.Errorf("get config for database '%s': %w", dbalias, err)
			}

			// if dbCfg.Engine == "postgres" {
			// postgres.Backup(dbCfg, backupDirPath, dbalias)
			// }

			return postgres.Backup(dbCfg, backupDirPath, dbalias)
		},
	}
	return backupCmd
}

func init() {
	// rootCmd.AddCommand(backupCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// backupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// backupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
