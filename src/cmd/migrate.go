package cmd

import (
	"assignment/datastore/mysql/migration"
	"assignment/logger"
	zaplogger "assignment/logger/zap"
	"github.com/spf13/cobra"
)

var forceMigrate bool = false

var MigrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate Base Project Database",
	Run: func(cmd *cobra.Command, args []string) {
		// Init Logger
		logger.Logger = zaplogger.NewLogger()

		initComponent()

		initTimezone()

		migration.Migrate(false, -1, forceMigrate, false)

		logger.SyncLogger()
	},
}

func init() {
	rootCmd.AddCommand(MigrateCmd)
	MigrateCmd.PersistentFlags().BoolVar(&forceMigrate, "force", false, "force migrate (default is false)")
}
