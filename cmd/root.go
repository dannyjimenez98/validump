package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dannyjimenez98/validump/internal/config"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type App struct {
	ViperConfig *viper.Viper
}

var cfgPath string
var rootCmd = newRootCmd()

func newRootCmd() *cobra.Command {
	app := &App{}
	root := &cobra.Command{
		Use:   "validump",
		Short: "A CLI tool for database backups and recovery validation.",
		Long: `Validump is a CLI-based database backup utility and disaster recovery validation tool.

		It allows for scheduled backups of your database and tests the backed up data in a containerized
		enviroment to ensure that the database can be successfully recovered in the event of an emergency.`,
		SilenceUsage: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if cmd.Annotations["skipLoadConfig"] == "true" {
				return nil
			}

			path, err := getConfigPath()
			if err != nil {
				return fmt.Errorf("get config path: %w", err)
			}

			if err := config.Validate(path); err != nil {
				return fmt.Errorf("validate config: %w", err)
			}

			v, err := config.Load(path)
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}

			app.ViperConfig = v

			return nil
		},
	}

	root.AddCommand(newPingCmd(app))
	root.SetErrPrefix(color.RedString("Error:"))
	root.PersistentFlags().StringVar(
		&cfgPath,
		"config",
		"",
		"Path to validump configuration file",
	)

	return root
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

	if envCfg := os.Getenv("VALIDUMP_CONFIG"); envCfg != "" {
		return filepath.Abs(envCfg)
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
