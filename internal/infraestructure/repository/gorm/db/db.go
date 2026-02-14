package db

import (
	"github.com/danielalmeidafarias/go_stock_engine/internal/domain"
	"gorm.io/gorm"
)

type GoORM interface {
	GetORM() *gorm.DB
	MapGormError(err error, context string) *domain.Error
}
