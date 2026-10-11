package cmd

import (
	"fmt"
	"os"

	"github.com/dannyjimenez98/validump/internal/config"
	"github.com/spf13/cobra"
	"go.yaml.in/yaml/v4"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use: "init",
	Short: "Create initial configuration file and database backup directory.",
	Long: `Create validump.yaml and a backups directory in the current directory.

	The configuration specifies the identifying name for the database (initialized as "default"),
	as well as the type of database (postgresql by default), the database connection string (reads
	from DB_CONNECTION_STR in .env file, and the filepath for the database backup directory
	(default: "./backups").
	
	Edit validump.yaml to match your database settings. Existing configuration
	files are never overwritten.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := &config.Config{
			Databases: map[string]config.DBConfig{
				"default": {
					Type: "postgres",
					ConnectionStr: "DB_CONNECTION_STR",
					BackupDir: "backups",
				},
			},
		}

		dbConfig, err := yaml.Marshal(&cfg)
		if err != nil {
			return fmt.Errorf("marshal database configuration: %w", err)
		}
		
		// create backups directory if it does not exist
		if err := os.Mkdir(cfg.Databases["default"].BackupDir, 0700); err != nil {
			return fmt.Errorf("create backup directory: %w", err)
		}

		filename := "validump.yaml"

		// create initial yaml config file if it does not exist
		file, err := os.OpenFile(filename, os.O_WRONLY|os.O_EXCL|os.O_CREATE, 0644)
		if err != nil {
			return fmt.Errorf("initialize validump configuration: %w", err)
		}
		defer file.Close()

		if _, err := file.Write(dbConfig); err != nil {
			return fmt.Errorf("write yaml configuration: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
