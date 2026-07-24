package main

import (
	"errors"
	"flag"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"

	"github.com/danielalmeidafarias/go_stock_engine/internal/config"
	"github.com/danielalmeidafarias/go_stock_engine/internal/infrastructure/migration"
)

const postgresDriver = "POSTGRES"

func main() {
	baselineVersion := flag.Int("baseline-version", -1, "mark an existing schema as the specified migration version")
	flag.Parse()

	database := config.Load().Database()
	if database.Driver != postgresDriver {
		log.Fatalf("unsupported migration driver: %s", database.Driver)
	}

	source, err := iofs.New(migration.Postgres, "postgres")
	if err != nil {
		log.Fatal("failed to load migrations: ", err)
	}

	migrator, err := migrate.NewWithSourceInstance("iofs", source, database.ConnectionString)
	if err != nil {
		log.Fatal("failed to initialize migrations: ", err)
	}
	defer migrator.Close()

	if *baselineVersion >= 0 {
		if err := migrator.Force(*baselineVersion); err != nil {
			log.Fatal("failed to set migration baseline: ", err)
		}
		log.Printf("migration baseline set to version %d", *baselineVersion)
		return
	}

	if err := migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal("failed to apply migrations: ", err)
	}
}
