package main

import (
	"errors"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"

	"github.com/danielalmeidafarias/go_stock_engine/internal/config"
	"github.com/danielalmeidafarias/go_stock_engine/internal/infrastructure/migration"
)

const postgresDriver = "POSTGRES"

func main() {
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

	if err := migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal("failed to apply migrations: ", err)
	}
}
