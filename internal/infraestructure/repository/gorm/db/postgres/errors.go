package postgres

import (
	"strings"

	"github.com/danielalmeidafarias/go_stock_engine/internal/domain"
	"github.com/jackc/pgx/v5/pgconn"
)

// Códigos de erro do PostgreSQL
// https://www.postgresql.org/docs/current/errcodes-appendix.html
const (
	pgUniqueViolation           = "23505"
	pgForeignKeyViolation       = "23503"
	pgNotNullViolation          = "23502"
	pgCheckViolation            = "23514"
	pgStringDataTooLong         = "22001"
	pgInvalidTextRepresentation = "22P02"
)

func (pgConn *PostgresConnection) MapGormError(pgErr *pgconn.PgError, context string) *domain.Error {
	switch pgErr.Code {
	case pgUniqueViolation:
		field := "field"
		parts := strings.Split(pgErr.ConstraintName, "_")
		if len(parts) >= 2 {
			field = parts[len(parts)-2]
		}

		return domain.NewError(context+": "+field+" already in use", domain.ErrConflict)

	case pgForeignKeyViolation:
		return domain.NewError(context+": invalid reference to related record", domain.ErrBadRequest)

	case pgNotNullViolation:
		return domain.NewError(context+": field '"+pgErr.ColumnName+"' is required", domain.ErrBadRequest)

	case pgCheckViolation:
		return domain.NewError(context+": invalid value for constraint '"+pgErr.ConstraintName+"'", domain.ErrBadRequest)

	case pgStringDataTooLong:
		return domain.NewError(context+": value exceeds maximum allowed length", domain.ErrBadRequest)

	case pgInvalidTextRepresentation:
		return domain.NewError(context+": invalid data format (e.g. malformed UUID)", domain.ErrBadRequest)

	default:
		return domain.NewError(context+": database error ("+pgErr.Code+")", domain.ErrInternal)
	}
}
