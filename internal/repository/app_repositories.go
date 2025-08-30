package repository

import (
	"database/sql"
	"ticket-api/internal/db/version"
)

type AppRepositories struct {
	Version *VersionRepository
}

func NewRepositories(db *sql.DB) *AppRepositories {
	return &AppRepositories{
		Version: &VersionRepository{queries: version.New(db)},
	}
}
