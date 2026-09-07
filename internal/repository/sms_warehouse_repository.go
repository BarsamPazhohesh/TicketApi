package repository

import (
	"context"
	"ticket-api/internal/db/sms_warehouse"
	"ticket-api/internal/dto"
	"ticket-api/internal/errx"
)

type SMSWarehouseRepository struct {
	queries *sms_warehouse.Queries
}

func NewSMSWarehouseRepository(queries *sms_warehouse.Queries) *SMSWarehouseRepository {
	return &SMSWarehouseRepository{
		queries: queries,
	}
}

func (repo *SMSWarehouseRepository) Create(ctx context.Context, smsTypeID int64, phone, message string, status dto.SMSStatus) (*dto.SMSWarehouseDTO, *errx.APIError) {
	record, err := repo.queries.CreateSMSWarehouseRecord(ctx, sms_warehouse.CreateSMSWarehouseRecordParams{
		SmsTypeID:           smsTypeID,
		ReceiverPhoneNumber: phone,
		Message:             message,
		Status:              int64(status),
	})
	if err != nil {
		return nil, errx.Respond(errx.ErrInternalServerError, err)
	}

	return dto.ToSMSWarehouseDTO(record), nil
}

func (repo *SMSWarehouseRepository) UpdateStatus(ctx context.Context, id int64, status dto.SMSStatus) *errx.APIError {
	err := repo.queries.UpdateSMSWarehouseStatus(ctx, sms_warehouse.UpdateSMSWarehouseStatusParams{
		ID:     id,
		Status: int64(status),
	})
	if err != nil {
		return errx.Respond(errx.ErrInternalServerError, err)
	}
	return nil
}

func (repo *SMSWarehouseRepository) GetPending(ctx context.Context) ([]dto.SMSWarehouseDTO, *errx.APIError) {
	records, err := repo.queries.GetPendingSMSWarehouseRecords(ctx)
	if err != nil {
		return nil, errx.Respond(errx.ErrInternalServerError, err)
	}

	result := make([]dto.SMSWarehouseDTO, 0, len(records))
	for _, r := range records {
		result = append(result, *dto.FromGetPendingSMSWarehouseRecordsRow(r))
	}
	return result, nil
}

func (repo *SMSWarehouseRepository) GetByID(ctx context.Context, id int64) (*dto.SMSWarehouseDTO, *errx.APIError) {
	record, err := repo.queries.GetSMSWarehouseByID(ctx, id)
	if err != nil {
		return nil, errx.Respond(errx.ErrInternalServerError, err)
	}
	return dto.FromGetSMSWarehouseByIDRow(record), nil
}
