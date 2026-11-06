package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/dannyjimenez98/validump/internal/config"
	"github.com/spf13/cobra"
	"go.yaml.in/yaml/v4"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use: "init",
	Short: "Create initial configuration file",
	Long: `Create validump.yaml file for the database configuration.

	The configuration file is generated with default placeholder values where the user can
	replace them with ones for their own database. By default, the file will be created at 
	~/.config/validump/validump.yaml (Mac, Linux) or %APPDATA%\validump\validump.yaml (Windows).

	The configuration specifies the alias for the database (initially set as "db-alias"),
	as well as the database engine, connection variables, and the filepath specifying where
	the directory that will contain the backup files for each database will be located
	(default: "~/validump/backups").
	
	Edit validump.yaml to match your database settings. Existing configuration
	files are never overwritten.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		filename := "validump.yaml"
		if _, err := os.Stat(filepath.Join(configPath, filename)); !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("%s already exists and was not overwritten", filepath.Join(configPath, filename))
		}

		// generate initial config template with default values
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("resolve home directory path: %w", err)
		}
		cfg := &config.Config{
			BackupDir: filepath.Join(home, "validump", "backups"),
			Databases: map[string]config.DBConfig{
				"db-alias": {
					Engine: "postgres",
					Host: "localhost",
					Port: 5432,
					Username: "dbuser",
					DBName: "dbname",
				},
			},
		}

		dbConfig, err := yaml.Marshal(&cfg)
		if err != nil {
			return fmt.Errorf("marshal database configuration: %w", err)
		}
		
		// create config directory path if it does not exist
		if err := os.MkdirAll(configPath, 0700); err != nil {
			return fmt.Errorf("create validump config directory: %w", err)
		}

		// create initial yaml config file if it does not exist
		file, err := os.OpenFile(filepath.Join(configPath, filename), os.O_WRONLY|os.O_EXCL|os.O_CREATE, 0644)
		if err != nil {
			return fmt.Errorf("create validump configuration template: %w", err)
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
