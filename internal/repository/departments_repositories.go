package repository

import (
	"context"
	"ticket-api/internal/db/departments"
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
	departmentId, err := repo.queries.AddDepartment(ctx, department)
	if err != nil {
		return -1, err
	}

	return departmentId, nil
}

func (repo *DepartmentsRepository) GetAllDepartments(ctx context.Context) ([]departments.Department, error) {
	return repo.queries.GetAllDepartments(ctx)
}
