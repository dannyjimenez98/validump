package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dannyjimenez98/validump/internal/config"
	"github.com/spf13/cobra"
)

var cfgPath string

var rootCmd = &cobra.Command{
	Use:   "validump",
	Short: "A CLI tool for database backups and recovery validation.",
	Long: `Validump is a CLI-based database backup utility and disaster recovery validation tool.

	It allows for scheduled backups of your database and tests the backed up data in a containerized
	enviroment to ensure that the database can be successfully recovered in the event of an emergency.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if cmd.Annotations["skipLoadConfig"] != "true" {
			path, err := getConfigPath()
			if err != nil {
				return fmt.Errorf("get config path: %w", err)
			}

			if _, err := config.Load(path); err != nil {
				return fmt.Errorf("load config: %w", err)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(
		&cfgPath,
		"config",
		"",
		"Path to validump configuration file",
	)
}

func Execute() error {
	return rootCmd.Execute()
}

func getConfigPath() (string, error) {
	if cfgPath != "" {
		// config paths explicitly set without a file extension
		// should have .yaml extension appended by default
		if filepath.Ext(cfgPath) == "" {
			cfgPath += ".yaml"
		}
		return filepath.Abs(cfgPath)
	}

	var cfgDir string
	var err error

	// check if $XDG_CONFIG_HOME is overriding default config path
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" && filepath.IsAbs(xdg) {
		cfgDir = xdg
	} else {
		// fallback to default user profile config path
		cfgDir, err = os.UserConfigDir()
		if err != nil {
			return "", fmt.Errorf("retrieve default configuration directory: %w", err)
		}
	}

	return filepath.Join(cfgDir, "validump", "validump.yaml"), nil
}
