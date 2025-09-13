package dto

import (
	"database/sql"
	"ticket-api/internal/db/departments"
	//"ticket-api/internal/db/departments"
)

type DepartmentDto struct {
	ID          int64          `json:"id"`
	Title       string         `json:"title"`
	Description sql.NullString `json:"description"`
}

func (dt *DepartmentDto) ToModel() *departments.Department {
	status := int64(1)
	deleted := int64(0)

	return &departments.Department{
		ID:          dt.ID,
		Title:       dt.Title,
		Description: sql.NullString{String: dt.Description.String, Valid: dt.Description.String != ""},
		Status:      status,
		Deleted:     deleted,
	}
}

type DepartmentWithStatusDto struct {
	ID          int64          `json:"id"`
	Title       string         `json:"title"`
	Description sql.NullString `json:"description"`
	Status      int64          `json:"status"`
}
