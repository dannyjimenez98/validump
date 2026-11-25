package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dannyjimenez98/validump/internal/backup"
	"github.com/dannyjimenez98/validump/internal/config"
	"github.com/spf13/cobra"
)

// newBackupCmd represents the backup command
func newBackupCmd(cfg *App) *cobra.Command {
	opts := backup.Options{}
	var backupCmd = &cobra.Command{
		Use:   "backup <db-alias>",
		Short: "Backup a database",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			dbalias := args[0]

			backupDir := cfg.ViperConfig.GetString("backup_dir")
			if backupDir == "" {
				home, err := os.UserHomeDir()
				if err != nil {
					return fmt.Errorf("get user home directory: %w", err)
				}
				backupDir = filepath.Join(home, "validump", dbalias)
			}

			if err := os.MkdirAll(backupDir, 0700); err != nil {
				return fmt.Errorf("create backups directory for database '%s': %w", dbalias, err)
			}

			dbCfg, err := config.LoadDB(cfg.ViperConfig, dbalias)
			if err != nil {
				return fmt.Errorf("get config for database '%s': %w", dbalias, err)
			}

			svc := backup.NewService(dbalias, dbCfg, backupDir, opts)

			return svc.Backup(cmd.Context())
		},
	}
	backupCmd.Flags().StringVarP(&opts.Format, "format", "F", "custom", "backup format: plain, custom, directory, tar (default: custom)")
	backupCmd.Flags().BoolVarP(&opts.Compression, "compress", "c", false, "compression type")
	backupCmd.Flags().DurationVarP(&opts.Timeout, "timeout", "t", 0, "maximum duration before timeout, e.g. 30s, 10m, 1h30m (default: 0 - timeout disabled)")
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
