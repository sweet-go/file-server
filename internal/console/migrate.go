// Package console hold all CLI-based interface functionality
package console

import (
	"strconv"

	migrate "github.com/rubenv/sql-migrate"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/sweet-go/file-server/internal/config"
	"github.com/sweet-go/file-server/internal/db"
	"github.com/sweet-go/stdlib/helper"
)

var runMigrate = &cobra.Command{
	Use:   "migrate",
	Short: "run database migration",
	Long:  "Used to run databse migration defined in migration folder",
	Run:   migration,
}

func init() {
	runMigrate.PersistentFlags().Int("step", 0, "maximum migration steps")
	runMigrate.PersistentFlags().String("direction", "up", "migration direction")
	RootCMD.AddCommand(runMigrate)
}

func migration(cmd *cobra.Command, _ []string) {
	direction := cmd.Flag("direction").Value.String()
	stepStr := cmd.Flag("step").Value.String()
	step, err := strconv.Atoi(stepStr)
	if err != nil {
		logrus.WithField("stepStr", stepStr).Fatal("Failed to parse step to int: ", err)
	}

	migrations := &migrate.FileMigrationSource{
		Dir: "./assets/migrations",
	}

	migrate.SetTable("schema_migrations")

	db.InitializePostgresConn()

	pgdb, err := db.PostgresDB.DB()
	if err != nil {
		logrus.WithField("DatabaseDSN", config.PostgresDSN()).Fatal("failed to run migration")
	}

	var n int
	if direction == "down" {
		n, err = migrate.ExecMax(pgdb, "postgres", migrations, migrate.Down, step)
	} else {
		n, err = migrate.ExecMax(pgdb, "postgres", migrations, migrate.Up, step)
	}
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"migrations": helper.Dump(migrations),
			"direction":  direction}).
			Fatal("Failed to migrate database: ", err)
	}

	logrus.Infof("Applied %d migrations!\n", n)
}
