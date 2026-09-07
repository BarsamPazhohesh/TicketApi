package dto

import (
	"database/sql"
	"ticket-api/internal/db/sms_warehouse"
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
	SMSTypeID           int64     `json:"smsTypeId"`
	SMSTypeTitle        string    `json:"smsTypeTitle,omitempty"`
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
		SMSTypeID:           m.SmsTypeID,
		ReceiverPhoneNumber: m.ReceiverPhoneNumber,
		Message:             m.Message,
		Status:              SMSStatus(m.Status),
		CreatedAt:           m.CreatedAt,
		UpdatedAt:           m.UpdatedAt,
		DeletedAt:           deletedAt,
	}
}

func FromGetPendingSMSWarehouseRecordsRow(row sms_warehouse.GetPendingSMSWarehouseRecordsRow) *SMSWarehouseDTO {
	var deletedAt *string
	if row.DeletedAt.Valid {
		deletedAt = &row.DeletedAt.String
	}

	return &SMSWarehouseDTO{
		ID:                  row.ID,
		SMSTypeID:           row.SmsTypeID,
		SMSTypeTitle:        row.SmsTypeTitle,
		ReceiverPhoneNumber: row.ReceiverPhoneNumber,
		Message:             row.Message,
		Status:              SMSStatus(row.Status),
		CreatedAt:           row.CreatedAt,
		UpdatedAt:           row.UpdatedAt,
		DeletedAt:           deletedAt,
	}
}

func FromGetSMSWarehouseByIDRow(row sms_warehouse.GetSMSWarehouseByIDRow) *SMSWarehouseDTO {
	var deletedAt *string
	if row.DeletedAt.Valid {
		deletedAt = &row.DeletedAt.String
	}

	return &SMSWarehouseDTO{
		ID:                  row.ID,
		SMSTypeID:           row.SmsTypeID,
		SMSTypeTitle:        row.SmsTypeTitle,
		ReceiverPhoneNumber: row.ReceiverPhoneNumber,
		Message:             row.Message,
		Status:              SMSStatus(row.Status),
		CreatedAt:           row.CreatedAt,
		UpdatedAt:           row.UpdatedAt,
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
		SmsTypeID:           dto.SMSTypeID,
		ReceiverPhoneNumber: dto.ReceiverPhoneNumber,
		Message:             dto.Message,
		Status:              int64(dto.Status),
		CreatedAt:           dto.CreatedAt,
		UpdatedAt:           dto.UpdatedAt,
		DeletedAt:           deletedAt,
	}
}

