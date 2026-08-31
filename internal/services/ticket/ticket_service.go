package ticket

import (
	"context"
	"errors"
	"fmt"
	"ticket-api/internal/config"
	"ticket-api/internal/dto"
	"ticket-api/internal/errx"
	"ticket-api/internal/model"
	"ticket-api/internal/repository"
	"ticket-api/internal/services/storage"
	"ticket-api/internal/util"
	"time"
)

type TicketService struct {
	ticketRepo         *repository.TicketRepository
	chatRepo           *repository.ChatRepository
	ticketTypeRepo     *repository.TicketTypesRepository
	ticketPriorityRepo *repository.TicketPrioritiesRepository
	ticketStatusRepo   *repository.TicketStatusesRepository
	userRepo           *repository.UsersRepository
	departmentRepo     *repository.DepartmentsRepository
	storageService     *storage.StorageService
}

func NewTicketService(
	ticketRepo *repository.TicketRepository,
	chatRepo *repository.ChatRepository,
	ticketTypeRepo *repository.TicketTypesRepository,
	ticketPriorityRepo *repository.TicketPrioritiesRepository,
	ticketStatusRepo *repository.TicketStatusesRepository,
	userRepo *repository.UsersRepository,
	departmentRepo *repository.DepartmentsRepository,
	storageService *storage.StorageService,
) *TicketService {
	return &TicketService{
		ticketRepo:         ticketRepo,
		chatRepo:           chatRepo,
		ticketTypeRepo:     ticketTypeRepo,
		ticketPriorityRepo: ticketPriorityRepo,
		ticketStatusRepo:   ticketStatusRepo,
		userRepo:           userRepo,
		departmentRepo:     departmentRepo,
		storageService:     storageService,
	}
}

// CreateTicket orchestrates ticket validation, track code generation, attachment promotion, and persistence
func (s *TicketService) CreateTicket(ctx context.Context, currentUserID int64, req dto.TicketCreateRequest) (*dto.TicketCreateResponse, *errx.APIError) {
	// 1. Validate attachment limit
	if len(req.Attachments) > 0 {
		total := len(req.Attachments)
		if total > config.Get().TicketConfig.MaxTicketUploadFile {
			apiErr := errx.Respond(errx.ErrMaxTicketFilesExceeded, errors.New(""))
			apiErr.Err.Message += fmt.Sprintf(" حداکثر فایل مجاز برای هر تیکت: %d", config.Get().TicketConfig.MaxTicketUploadFile)
			return nil, apiErr
		}
	}

	// 2. Set authenticated user ID (prevent spoofing if currentUserID > 0)
	targetUserID := req.UserID
	if currentUserID > 0 {
		targetUserID = currentUserID
	}

	// 3. Check user exists
	isUserExist, apiErr := s.userRepo.IsUserExist(ctx, targetUserID)
	if apiErr != nil {
		return nil, apiErr
	}
	if !isUserExist {
		return nil, errx.Respond(errx.ErrUserNotFound, errors.New("user not found"))
	}

	// 4. Check ticket type exists
	isTicketTypeExists, apiErr := s.ticketTypeRepo.IsTicketTypeExits(ctx, req.TicketTypeID)
	if apiErr != nil {
		return nil, apiErr
	}
	if !isTicketTypeExists {
		return nil, errx.Respond(errx.ErrTicketTypeNotFound, errors.New("ticket type not found"))
	}

	// 5. Check department exists
	isDepExists, apiErr := s.departmentRepo.IsDepartmentExits(ctx, req.DepartmentID)
	if apiErr != nil {
		return nil, apiErr
	}
	if !isDepExists {
		return nil, errx.Respond(errx.ErrDepartmentNotFound, errors.New("department not found"))
	}

	// 6. Get open status
	openStatus, apiErr := s.ticketStatusRepo.GetOpenStatus(ctx)
	if apiErr != nil {
		return nil, apiErr
	}

	// 7. Parse and clean attachment filenames
	attachments, err := util.ParseObjectNames(req.Attachments)
	if err != nil {
		return nil, errx.Respond(errx.ErrBadRequest, err)
	}

	now := time.Now()
	ticketID := util.GenerateUUID()
	trackCode, err := util.GenerateUniqueTrackCode(ctx, s.ticketRepo.GetCollection())
	if err != nil {
		return nil, errx.Respond(errx.ErrInternalServerError, err)
	}

	// 8. Move temp files to ticket directory if present
	movedAttachments := attachments
	if len(attachments) > 0 && s.storageService != nil {
		moved, moveErr := s.storageService.MoveTempsFileToTickets(ctx, ticketID, attachments)
		if moveErr != nil {
			return nil, moveErr
		}
		movedAttachments = moved
	}

	firstMessage := model.ChatMessage{
		ID:          util.GenerateUUID(),
		SenderID:    targetUserID,
		Message:     req.Body,
		Attachments: movedAttachments,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	ticketModel := &model.Ticket{
		ID:              ticketID,
		TrackCode:       trackCode,
		UserID:          targetUserID,
		TicketTypeID:    req.TicketTypeID,
		DepartmentID:    req.DepartmentID,
		TicketStatusID:  openStatus.ID,
		Title:           req.Title,
		AttachmentCount: len(movedAttachments),
		CreatedAt:       now,
		UpdatedAt:       now,
		Chat:            []model.ChatMessage{firstMessage},
	}

	if err := s.ticketRepo.InsertTicket(ctx, ticketModel); err != nil {
		return nil, err
	}

	return &dto.TicketCreateResponse{
		ID:        ticketID,
		TrackCode: trackCode,
	}, nil
}

// CreateChatMessage adds a message to ticket and moves attachments
func (s *TicketService) CreateChatMessage(ctx context.Context, ticketID string, senderID int64, req dto.ChatMessageCreateRequest) (*dto.ChatMessageDTO, *errx.APIError) {
	if len(req.Attachments) > 0 {
		currentCount, apiErr := s.ticketRepo.GetTicketAttachmentCount(ctx, ticketID)
		if apiErr != nil {
			return nil, apiErr
		}
		total := currentCount + len(req.Attachments)
		if total > config.Get().TicketConfig.MaxTicketUploadFile {
			apiErr := errx.Respond(errx.ErrMaxTicketFilesExceeded, errors.New(""))
			apiErr.Err.Message += fmt.Sprintf(" حداکثر فایل مجاز برای هر تیکت: %d", config.Get().TicketConfig.MaxTicketUploadFile)
			return nil, apiErr
		}
	}

	attachments, err := util.ParseObjectNames(req.Attachments)
	if err != nil {
		return nil, errx.Respond(errx.ErrBadRequest, err)
	}

	movedAttachments := attachments
	if len(attachments) > 0 && s.storageService != nil {
		moved, moveErr := s.storageService.MoveTempsFileToTickets(ctx, ticketID, attachments)
		if moveErr != nil {
			return nil, moveErr
		}
		movedAttachments = moved
	}

	now := time.Now()
	msgModel := model.ChatMessage{
		ID:          util.GenerateUUID(),
		SenderID:    senderID,
		Message:     req.Message,
		Attachments: movedAttachments,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	return s.chatRepo.AppendChatMessage(ctx, ticketID, msgModel)
}

// GetTicketByID fetches ticket by ID with IDOR authorization protection
func (s *TicketService) GetTicketByID(ctx context.Context, ticketID string, currentUserID int64, roleIDs []int64) (*dto.TicketResponse, *errx.APIError) {
	ticket, apiErr := s.ticketRepo.GetTicketByID(ctx, ticketID)
	if apiErr != nil {
		return nil, apiErr
	}

	// IDOR check: allow if admin/agent (roles 1 or 2) or if owner
	isAdminOrAgent := false
	for _, r := range roleIDs {
		if r == 1 || r == 2 {
			isAdminOrAgent = true
			break
		}
	}

	if !isAdminOrAgent && currentUserID > 0 && ticket.UserID != currentUserID {
		return nil, errx.Respond(errx.ErrForbidden, errors.New("access denied to this ticket"))
	}

	return ticket, nil
}

// GetTicketByTrackCode fetches ticket verifying username ownership
func (s *TicketService) GetTicketByTrackCode(ctx context.Context, trackCode, username string) (*dto.TicketResponse, *errx.APIError) {
	_, parseErr := util.ParsTrackCode(trackCode)
	if parseErr != nil {
		return nil, errx.Respond(errx.ErrBadRequest, parseErr)
	}

	user, apiErr := s.userRepo.GetUserByUsername(ctx, username)
	if apiErr != nil {
		if apiErr.Err.Code == errx.ErrUserNotFound {
			return nil, errx.Respond(errx.ErrTicketNotFound, errors.New("username not found"))
		}
		return nil, apiErr
	}

	ticket, apiErr := s.ticketRepo.GetTicketByTrackCode(ctx, trackCode)
	if apiErr != nil {
		return nil, apiErr
	}

	if ticket.UserID != user.ID {
		return nil, errx.Respond(errx.ErrTicketNotFound, errors.New("this username did not create this ticket"))
	}

	return ticket, nil
}

func (s *TicketService) GetTicketsList(ctx context.Context, params dto.TicketQueryParams) (*dto.TicketPagingResponse, *errx.APIError) {
	return s.ticketRepo.GetTickets(ctx, params)
}

func (s *TicketService) GetAllActiveTicketTypes(ctx context.Context) ([]*dto.TicketTypeDto, *errx.APIError) {
	ticketTypesList, err := s.ticketTypeRepo.GetAllActiveTicketTypes(ctx)
	if err != nil {
		return nil, err
	}

	ticketTypesDTO := make([]*dto.TicketTypeDto, 0, len(ticketTypesList))
	for _, v := range ticketTypesList {
		ticketTypesDTO = append(ticketTypesDTO, dto.ToTicketTypeDTO(&v))
	}
	return ticketTypesDTO, nil
}

func (s *TicketService) GetAllActiveTicketStatuses(ctx context.Context) ([]dto.TicketStatusDTO, *errx.APIError) {
	ticketStatuses, err := s.ticketStatusRepo.GetAllActiveTicketStatuses(ctx)
	if err != nil {
		return nil, err
	}

	ticketStatusDTO := make([]dto.TicketStatusDTO, 0, len(ticketStatuses))
	for _, v := range ticketStatuses {
		ticketStatusDTO = append(ticketStatusDTO, *dto.ToTicketStatusDTO(&v))
	}
	return ticketStatusDTO, nil
}

func (s *TicketService) CloseTicket(ctx context.Context, ticketID string) (*dto.TicketResponse, *errx.APIError) {
	closeStatus, err := s.ticketStatusRepo.GetCloseStatus(ctx)
	if err != nil {
		return nil, err
	}
	return s.ticketRepo.SetTicketStatus(ctx, ticketID, closeStatus.ID)
}
