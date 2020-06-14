package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var upgradeCommand = &cobra.Command{
	Use:   "upgrade",
	Short: "Update the DB schema to the newest version",
	Long: "Applies all migration files one by one according to the timestamp index. The newest vers" +
		"ion timestamp is stored in the specific DB table. Each DB upgrade executes all migrations " +
		"from the first one non-applied to the last one non-applied. After production usage migrati" +
		"on scripts can't be removed or modified `cause it breaks the versioning order.",
	Run: func(command *cobra.Command, args []string) {
		log.Debug("cmd: upgraded hello world!")
	},
}
