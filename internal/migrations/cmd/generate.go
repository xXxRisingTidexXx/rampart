package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var generateCommand = &cobra.Command{
	Use:   "generate",
	Short: "Create a new migration",
	Long: "Generate a new SQL-migration with the specific timestamp-base name. As soon as all times" +
		"tamps go in the ascending order, these file names are in the same time indices helping dis" +
		"tinct version upgrade sequence.",
	Run: func(command *cobra.Command, args []string) {
		log.Debug("cmd: generated hello world!")
	},
}
