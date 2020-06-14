package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var generateCommand = &cobra.Command{
	Use:   "generate",
	Short: "Create a new migration",
	Long: "Generate a new SQL-migration with the specific timestamp-base name. As\n" +
		"soon as all timestamps go in the ascending order, these file names are\n" +
		"in the same time indices helping distinct version upgrade sequence.",
	Run: func(command *cobra.Command, args []string) {
		log.Debug("cmd: generate started")
		log.Debug("cmd: generate finished")
	},
}
