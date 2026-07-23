package factory

import (
	"github.com/danielalmeidafarias/go_stock_engine/internal/domain/repository"
	"github.com/danielalmeidafarias/go_stock_engine/internal/infrastructure/repository/db"
	"github.com/danielalmeidafarias/go_stock_engine/internal/infrastructure/repository/db/postgres"
)

type Type string

const Postgres Type = "POSTGRES"

func NewProductStockRepository(repositoryType Type) repository.IProductStockRepository {
	switch repositoryType {
	case Postgres:
		return db.NewProductStockRepository(postgres.NewPostgresConnection())
	default:
		panic("invalid database type")
	}
}
