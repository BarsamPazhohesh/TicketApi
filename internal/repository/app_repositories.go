package repository

import (
	"database/sql"
	"ticket-api/internal/db/version"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type AppRepositories struct {
	Version *VersionRepository
	Ticket  *TicketRepository
}

func NewRepositories(sqldb *sql.DB, mongodb *mongo.Database) *AppRepositories {
	return &AppRepositories{
		Version: NewVersionRepository(version.New(sqldb)),
		Ticket:  NewTicketRepository(mongodb),
	}
}
