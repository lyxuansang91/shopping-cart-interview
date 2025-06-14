package main

import (
	"log"
	"os"

	"github.com/cinchprotocol/cinch-api/packages/dbtool/internal"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "dbtool",
		Short: "Database tool for running migrations and seeders across all services",
		Long:  `A CLI tool for managing database operations across all services in the Cinch API.`,
	}

	// Add commands
	rootCmd.AddCommand(internal.MigrateCmd)
	rootCmd.AddCommand(internal.SeedCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
