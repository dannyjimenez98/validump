package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"
)

var	configPath string

var rootCmd = &cobra.Command{
	Use: "validump",
	Short: "A CLI tool for database backups and recovery validation.",
	Long: `Validump is a CLI-based database backup utility and disaster recovery validation tool.

	It allows for scheduled backups of your database and tests the backed up data in a containerized
	enviroment to ensure that the database can be successfully recovered in the event of an emergency.`,
}

func init() {
	rootCmd.PersistentFlags().StringVar(
		&configPath,
		"config",
		getDefaultConfigPath(),
		"Path to validump configuration file",
	)
}

func Execute() error {
	return rootCmd.Execute(); 
}

func getDefaultConfigPath() string {
	configRoot := "validump"
	// check if XDG_CONFIG_HOME variable is overriding default config path to ~/.config
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, configRoot)
	}

	// use ~/.config path for Mac and Linux systems
	if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Println("retrieve user home directory: ", err)
			return ""
		}	
		return filepath.Join(homeDir, ".config", configRoot)
	}

	// fallback to default user profile config path (default path for Windows systems)
	defaultConfigDir, err := os.UserConfigDir()
	if err != nil {
		fmt.Println("retrieve default configuration directory: ", err)
		return ""
	} 

	configPath := filepath.Join(defaultConfigDir, configRoot)

	return configPath
}
