package factory

import (
	"github.com/danielalmeidafarias/go_stock_engine/internal/config"
	"github.com/danielalmeidafarias/go_stock_engine/internal/domain/repository"
	"github.com/danielalmeidafarias/go_stock_engine/internal/infrastructure/repository/db"
	"github.com/danielalmeidafarias/go_stock_engine/internal/infrastructure/repository/db/postgres"
)

type Driver string

const Postgres Driver = "POSTGRES"

func NewProductStockRepository(databaseConfig config.DatabaseConfig) repository.IProductStockRepository {
	switch Driver(databaseConfig.Driver) {
	case Postgres:
		return db.NewProductStockRepository(postgres.NewPostgresConnection(databaseConfig.ConnectionString))
	default:
		panic("invalid database type")
	}
}
