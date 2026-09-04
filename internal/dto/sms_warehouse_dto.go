package dto

import (
	"database/sql"
	"ticket-api/internal/model"
)

type SMSStatus int64

const (
	SMSStatusPending SMSStatus = 0
	SMSStatusSent    SMSStatus = 1
	SMSStatusFailed  SMSStatus = 2
)

type SMSWarehouseDTO struct {
	ID                  int64     `json:"id"`
	ReceiverPhoneNumber string    `json:"receiverPhoneNumber"`
	Message             string    `json:"message"`
	Status              SMSStatus `json:"status"`
	CreatedAt           string    `json:"createdAt"`
	UpdatedAt           string    `json:"updatedAt"`
	DeletedAt           *string   `json:"deletedAt,omitempty"`
}

func ToSMSWarehouseDTO(m model.SMSWarehouse) *SMSWarehouseDTO {
	var deletedAt *string
	if m.DeletedAt.Valid {
		deletedAt = &m.DeletedAt.String
	}

	return &SMSWarehouseDTO{
		ID:                  m.ID,
		ReceiverPhoneNumber: m.ReceiverPhoneNumber,
		Message:             m.Message,
		Status:              SMSStatus(m.Status),
		CreatedAt:           m.CreatedAt,
		UpdatedAt:           m.UpdatedAt,
		DeletedAt:           deletedAt,
	}
}

func (dto *SMSWarehouseDTO) ToModel() *model.SMSWarehouse {
	var deletedAt sql.NullString
	if dto.DeletedAt != nil {
		deletedAt = sql.NullString{String: *dto.DeletedAt, Valid: true}
	}

	return &model.SMSWarehouse{
		ID:                  dto.ID,
		ReceiverPhoneNumber: dto.ReceiverPhoneNumber,
		Message:             dto.Message,
		Status:              int64(dto.Status),
		CreatedAt:           dto.CreatedAt,
		UpdatedAt:           dto.UpdatedAt,
		DeletedAt:           deletedAt,
	}
}
