package repository

import (
	"context"
	"database/sql"
	"errors"
	"ticket-api/internal/config"
	"ticket-api/internal/db/ticket_types"
	"ticket-api/internal/errx"
	"ticket-api/internal/services/cache"
	"time"
)

const (
	_TicketTypesAllKey    = "ticket_types_all"
	_TicketTypesActiveKey = "ticket_types_active"
)

type TicketTypesRepository struct {
	queries *ticket_types.Queries
	cache   *cache.CacheService
}

// pass cache service in constructor
func NewTicketTypesRepository(queries *ticket_types.Queries, cache *cache.CacheService) *TicketTypesRepository {
	return &TicketTypesRepository{
		queries: queries,
		cache:   cache,
	}
}

func (repo *TicketTypesRepository) AddTicketType(ctx context.Context, ticketType ticket_types.AddTicketTypeParams) (int64, *errx.APIError) {
	ticketTypeID, err := repo.queries.AddTicketType(ctx, ticketType)
	if err != nil {
		return -1, errx.Respond(errx.ErrInternalServerError, err)
	}

	// invalidate cache when new type added
	_ = repo.cache.Delete(ctx, _TicketTypesAllKey)
	_ = repo.cache.Delete(ctx, _TicketTypesActiveKey)

	return ticketTypeID, nil
}

// read all ticket types with cache
func (repo *TicketTypesRepository) GetAllTicketTypes(ctx context.Context) ([]ticket_types.TicketType, *errx.APIError) {
	var types []ticket_types.TicketType

	ok, err := repo.cache.Get(ctx, _TicketTypesAllKey, &types)
	if err != nil {
		return nil, errx.Respond(errx.ErrInternalServerError, err)
	}
	if ok {
		return types, nil
	}

	types, err = repo.queries.GetAllTicketTypes(ctx)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errx.Respond(errx.ErrTicketTypeNotFound, err)
		}
		return nil, errx.Respond(errx.ErrInternalServerError, err)
	}

	_ = repo.cache.Set(ctx, _TicketTypesAllKey, types, time.Duration(config.Get().Cache.TicketTypeTTL)*time.Minute)
	return types, nil
}

// read all active ticket types with cache
func (repo *TicketTypesRepository) GetAllActiveTicketTypes(ctx context.Context) ([]ticket_types.TicketType, *errx.APIError) {
	var types []ticket_types.TicketType

	ok, err := repo.cache.Get(ctx, _TicketTypesActiveKey, &types)
	if err != nil {
		return nil, errx.Respond(errx.ErrInternalServerError, err)
	}
	if ok {
		return types, nil
	}

	data, err := repo.queries.GetAllActiveTicketTypes(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errx.Respond(errx.ErrTicketTypeNotFound, err)
		}
		return nil, errx.Respond(errx.ErrInternalServerError, err)
	}

	_ = repo.cache.Set(ctx, _TicketTypesActiveKey, data, time.Duration(config.Get().Cache.TicketTypeTTL)*time.Minute)
	return data, nil
}

func (repo *TicketTypesRepository) IsTicketTypeExits(ctx context.Context, typeID int64) (bool, *errx.APIError) {
	count, err := repo.queries.CheckTicketTypeByID(ctx, typeID)
	if err != nil {
		return false, errx.Respond(errx.ErrInternalServerError, err)
	}
	return count != 0, nil
}
