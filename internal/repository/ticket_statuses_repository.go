package repository

import (
	"context"
	"database/sql"
	"errors"
	"ticket-api/internal/config"
	"ticket-api/internal/db/ticket_statuses"
	"ticket-api/internal/errx"
	"ticket-api/internal/services/cache"
	"time"
)

type TicketStatusesRepository struct {
	queries *ticket_statuses.Queries
	cache   *cache.CacheService
}

// private cache keys
const (
	_ticketStatusesAllKey = "ticket_status_all"
	_ticketStatusCloseKey = "ticket_status_close"
	_ticketStatusOpenKey  = "ticket_status_open"
)

func NewTicketStatusesRepository(queries *ticket_statuses.Queries, cache *cache.CacheService) *TicketStatusesRepository {
	return &TicketStatusesRepository{
		queries: queries,
		cache:   cache,
	}
}

// Add a new ticket status and invalidate caches
func (repo *TicketStatusesRepository) AddTicketStatus(ctx context.Context, param ticket_statuses.AddTicketStatusParams) *errx.APIError {
	err := repo.queries.AddTicketStatus(ctx, param)
	if err != nil {
		return errx.Respond(errx.ErrInternalServerError, err)
	}

	// invalidate cache
	_ = repo.cache.Delete(ctx, _ticketStatusesAllKey)
	_ = repo.cache.Delete(ctx, _ticketStatusOpenKey)
	_ = repo.cache.Delete(ctx, _ticketStatusCloseKey)
	return nil
}

// Get all active ticket statuses with cache
func (repo *TicketStatusesRepository) GetAllActiveTicketStatuses(ctx context.Context) ([]ticket_statuses.TicketStatus, *errx.APIError) {
	var statuses []ticket_statuses.TicketStatus

	// check cache first
	ok, err := repo.cache.Get(ctx, _ticketStatusesAllKey, &statuses)
	if err != nil {
		return nil, errx.Respond(errx.ErrInternalServerError, err)
	}
	if ok {
		return statuses, nil
	}

	// fallback to DB
	statuses, err = repo.queries.GetAllActiveTicketStatuses(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errx.Respond(errx.ErrTicketStatusNotFound, err)
		}
		return nil, errx.Respond(errx.ErrInternalServerError, err)
	}

	// set caches
	ttl := time.Duration(config.Get().Cache.TicketStatusTTL) * time.Minute
	_ = repo.cache.Set(ctx, _ticketStatusesAllKey, statuses, ttl)
	repo.cacheCloseOpen(ctx, statuses)

	return statuses, nil
}

// Get a single active ticket status by ID
func (repo *TicketStatusesRepository) GetActiveTicketStatusByID(ctx context.Context, ID int64) (*ticket_statuses.TicketStatus, *errx.APIError) {
	statuses, err := repo.GetAllActiveTicketStatuses(ctx)
	if err != nil {
		return nil, err
	}

	for _, s := range statuses {
		if s.ID == ID {
			return &s, nil
		}
	}

	return nil, errx.Respond(errx.ErrTicketStatusNotFound, nil)
}

// Get the "close" ticket status
func (repo *TicketStatusesRepository) GetCloseStatus(ctx context.Context) (*ticket_statuses.TicketStatus, *errx.APIError) {
	return repo.getCachedStatus(ctx, _ticketStatusCloseKey, 0)
}

// Get the "open" ticket status
func (repo *TicketStatusesRepository) GetOpenStatus(ctx context.Context) (*ticket_statuses.TicketStatus, *errx.APIError) {
	return repo.getCachedStatus(ctx, _ticketStatusOpenKey, 1)
}

// --- internal helper to reduce repetition ---
func (repo *TicketStatusesRepository) getCachedStatus(ctx context.Context, key string, index int) (*ticket_statuses.TicketStatus, *errx.APIError) {
	var status ticket_statuses.TicketStatus

	ok, err := repo.cache.Get(ctx, key, &status)
	if err != nil {
		return nil, errx.Respond(errx.ErrInternalServerError, err)
	}
	if ok {
		return &status, nil
	}

	// fallback to DB
	statuses, apiErr := repo.GetAllActiveTicketStatuses(ctx)
	if apiErr != nil {
		return nil, errx.Respond(errx.ErrInternalServerError, err)
	}

	if len(statuses) <= index {
		return nil, errx.Respond(errx.ErrInternalServerError, errors.New("statuses list too short"))
	}

	status = statuses[index]

	// store in cache
	ttl := time.Duration(config.Get().Cache.TicketStatusTTL) * time.Minute
	_ = repo.cache.Set(ctx, key, status, ttl)

	return &status, nil
}

// --- helper to cache "close" and "open" separately ---
func (repo *TicketStatusesRepository) cacheCloseOpen(ctx context.Context, statuses []ticket_statuses.TicketStatus) {
	if len(statuses) >= 2 {
		ttl := time.Duration(config.Get().Cache.TicketStatusTTL) * time.Minute
		_ = repo.cache.Set(ctx, _ticketStatusCloseKey, statuses[0], ttl)
		_ = repo.cache.Set(ctx, _ticketStatusOpenKey, statuses[1], ttl)
	}
}
