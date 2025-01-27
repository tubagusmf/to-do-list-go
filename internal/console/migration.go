package console

import (
	"database/sql"
	"log"
	"to-do-list/internal/config"
	"to-do-list/internal/helper"

	_ "github.com/lib/pq"

	migrate "github.com/rubenv/sql-migrate"
	"github.com/spf13/cobra"
)

var (
	direction string
	step      int
)

func init() {
	rootCmd.AddCommand(migrationCMD)

	migrationCMD.Flags().StringVarP(&direction, "direction", "d", "up", "Migration direction (up or down)")
	migrationCMD.Flags().IntVarP(&step, "step", "s", 1, "Number of migrations to apply")
}

var migrationCMD = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	Long:  `This command is used to apply or rollback database migrations.`,
	Run:   migrateDB,
}

func migrateDB(cmd *cobra.Command, args []string) {
	config.LoadWithViper()

	// Open database connection
	connDB, err := sql.Open("postgres", helper.GetConnectionString()) // Menggunakan driver "postgres"
	log.Printf("Connection String: %s", helper.GetConnectionString())

	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer connDB.Close()

	// Define migration source
	migrations := &migrate.FileMigrationSource{Dir: "./db/migrations"}

	var n int
	if direction == "down" {
		n, err = migrate.ExecMax(connDB, "postgres", migrations, migrate.Down, step)
	} else {
		n, err = migrate.ExecMax(connDB, "postgres", migrations, migrate.Up, step)
	}

	if err != nil {
		log.Fatalf("Error applying migrations: %v", err)
	}

	log.Printf("Successfully applied %d migrations", n)
}
