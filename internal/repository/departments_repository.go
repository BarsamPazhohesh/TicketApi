package repository

import (
	"context"
	"ticket-api/internal/db/departments"
	"ticket-api/internal/errx"
)

type DepartmentsRepository struct {
	queries *departments.Queries
}

func NewDepartmentsRepository(queries *departments.Queries) *DepartmentsRepository {
	return &DepartmentsRepository{
		queries: queries,
	}
}

func (repo *DepartmentsRepository) AddDepartment(ctx context.Context, department departments.AddDepartmentParams) (int64, error) {
	departmentID, err := repo.queries.AddDepartment(ctx, department)
	if err != nil {
		return -1, err
	}

	return departmentID, nil
}

func (repo *DepartmentsRepository) GetAllDepartments(ctx context.Context) ([]departments.Department, error) {
	return repo.queries.GetAllDepartments(ctx)
}

func (repo *DepartmentsRepository) IsDepartmentExits(ctx context.Context, departmentID int64) (bool, *errx.APIError) {
	count, err := repo.queries.CheckDepartmentByID(ctx, departmentID)
	if err != nil {
		return false, errx.Respond(errx.ErrInternalServerError, err)
	}
	return count != 0, nil
}
