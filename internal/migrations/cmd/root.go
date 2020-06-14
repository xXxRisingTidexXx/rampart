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
	Short: "Migrations is the CLI for Rampart project DB migrations",
	Long: "Quick and dirty PostgreSQL migration executor. Capable of an automatic SQL-template gene" +
		"ration and forward DB schema upgrade.",
}

func Execute() {
	if err := rootCommand.Execute(); err != nil {
		log.Fatal(err)
	}
}
