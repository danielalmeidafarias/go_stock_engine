package migration

import "embed"

// Postgres contains the versioned PostgreSQL migrations applied by cmd/migrate.
//
//go:embed postgres/*.sql
var Postgres embed.FS
