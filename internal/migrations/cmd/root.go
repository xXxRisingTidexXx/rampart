package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

//nolint:gochecknoinits
func init() {
	rootCommand.AddCommand(generateCommand, upgradeCommand)
}

var rootCommand = &cobra.Command{
	Use:   "migrations",
	Short: "Rampart DB migration tool",
	Long:  "Quick and dirty PostgreSQL DB migration CLI.",
}

func Execute() {
	if err := rootCommand.Execute(); err != nil {
		log.Fatal(err)
	}
}
