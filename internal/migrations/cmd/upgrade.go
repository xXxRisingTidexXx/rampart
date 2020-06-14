package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var upgradeCommand = &cobra.Command{
	Use:   "upgrade",
	Short: "Update the DB schema to the newest version",
	Long: "Applies all migration files one by one according to the timestamp index.\n" +
		"The newest version timestamp is stored in the specific DB table. Each DB\n" +
		"upgrade executes all migrations from the first one non-applied to the\n" +
		"last one non-applied. After production usage migration scripts can't be\n" +
		"removed or modified `cause it breaks the versioning order.",
	Run: func(command *cobra.Command, args []string) {
		log.Debug("cmd: upgrade started")
		log.Debug("cmd: upgrade finished")
	},
}
