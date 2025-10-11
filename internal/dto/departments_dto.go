package dto

import (
	"database/sql"
	"ticket-api/internal/db/departments"
	"ticket-api/internal/model"
	//"ticket-api/internal/db/departments"
)

type DepartmentDTO struct {
	ID          int64   `json:"id"`
	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`
}

func (dt *DepartmentDTO) ToModel() *departments.Department {
	nullDesc := sql.NullString{}
	if dt.Description != nil {
		nullDesc = sql.NullString{String: *dt.Description, Valid: true}
	} else {
		nullDesc = sql.NullString{String: "", Valid: false}
	}

	return &departments.Department{
		ID:          dt.ID,
		Title:       dt.Title,
		Description: nullDesc,
		Status:      1,
		Deleted:     0,
	}
}

func ToDepartmentDTO(m model.Department) *DepartmentDTO {
	var description *string
	if m.Description.Valid {
		description = &m.Description.String
	}

	return &DepartmentDTO{
		ID:          m.ID,
		Title:       m.Title,
		Description: description,
	}
}

type DepartmentWithStatusDto struct {
	ID          int64          `json:"id"`
	Title       string         `json:"title"`
	Description sql.NullString `json:"description"`
	Status      int64          `json:"status"`
}
