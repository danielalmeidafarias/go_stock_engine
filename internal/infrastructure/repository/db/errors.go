package db

import (
	"errors"
	"fmt"

	"github.com/danielalmeidafarias/go_stock_engine/internal/domain"
	"gorm.io/gorm"
)

func (repo *ProductStockRepository) MapErrorToDomain(err error, context string) *domain.Error {
	message := fmt.Sprintf("%s: %s", context, err.Error())

	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return domain.NewError(message, domain.ErrBadRequest)

	case errors.Is(err, gorm.ErrDuplicatedKey):
		return domain.NewError(message, domain.ErrConflict)

	case errors.Is(err, gorm.ErrForeignKeyViolated):
		return domain.NewError(message, domain.ErrBadRequest)

	default:
		return domain.NewError(message, domain.ErrInternal)
	}
}
