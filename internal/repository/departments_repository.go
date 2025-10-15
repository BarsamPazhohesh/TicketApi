package repository

import (
	"context"
	"ticket-api/internal/config"
	"ticket-api/internal/db/departments"
	"ticket-api/internal/errx"
	"ticket-api/internal/services/cache"
	"time"
)

const (
	CacheKeyTicketTypesAll    = "ticket_types_all"
	CacheKeyTicketTypesActive = "ticket_types_active"
	CacheKeyDepartmentsAll    = "departments_all"
	CacheKeyTicketStatusAll   = "ticket_status_all"
)

type DepartmentsRepository struct {
	queries *departments.Queries
	cache   *cache.CacheService
}

func NewDepartmentsRepository(queries *departments.Queries, cache *cache.CacheService) *DepartmentsRepository {
	return &DepartmentsRepository{
		queries: queries,
		cache:   cache,
	}
}

func (repo *DepartmentsRepository) AddDepartment(ctx context.Context, department departments.AddDepartmentParams) (int64, *errx.APIError) {
	departmentID, err := repo.queries.AddDepartment(ctx, department)
	if err != nil {
		return -1, errx.Respond(errx.ErrInternalServerError, err)
	}

	// invalidate cache when new department added
	_ = repo.cache.Delete(ctx, "departments_all")

	return departmentID, nil
}

// get all departments with cache
func (repo *DepartmentsRepository) GetAllDepartments(ctx context.Context) ([]departments.Department, *errx.APIError) {
	var depts []departments.Department

	ok, err := repo.cache.Get(ctx, "departments_all", &depts)
	if err != nil {
		return nil, errx.Respond(errx.ErrInternalServerError, err)
	}
	if ok {
		return depts, nil
	}

	depts, err = repo.queries.GetAllDepartments(ctx)
	if err != nil {
		return nil, errx.Respond(errx.ErrInternalServerError, err)
	}

	// set cache with TTL from YAML
	_ = repo.cache.Set(ctx, "departments_all", depts, time.Duration(config.Get().Cache.DepartmentTTL)*time.Minute)
	return depts, nil
}

func (repo *DepartmentsRepository) IsDepartmentExits(ctx context.Context, departmentID int64) (bool, *errx.APIError) {
	count, err := repo.queries.CheckDepartmentByID(ctx, departmentID)
	if err != nil {
		return false, errx.Respond(errx.ErrInternalServerError, err)
	}
	return count != 0, nil
}
