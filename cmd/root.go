package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "validump",
	Short: "A CLI tool for database backups and recovery validation.",
	Long: `Validump is a CLI-based database backup utility and disaster recovery validation tool.

	It allows for scheduled backups of your database and tests the backed up data in a containerized
	enviroment to ensure that the database can be successfully recovered in the event of an emergency.`,
}


func Execute() error {
	return rootCmd.Execute(); 
}
