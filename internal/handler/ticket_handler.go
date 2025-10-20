package handler

import (
	"errors"
	"fmt"
	"net/http"
	"ticket-api/internal/config"
	"ticket-api/internal/dto"
	"ticket-api/internal/errx"
	"ticket-api/internal/repository"
	"ticket-api/internal/util"

	_ "ticket-api/internal/routes"

	"github.com/gin-gonic/gin"
)

// TicketHandler handles ticket-related HTTP requests
type TicketHandler struct {
	TicketRepo         *repository.TicketRepository
	TicketTypeRepo     *repository.TicketTypesRepository
	TicketPriorityRepo *repository.TicketPrioritiesRepository
	TicketStatusRepo   *repository.TicketStatusesRepository
	UserRepo           *repository.UsersRepository
	DepartmentRepo     *repository.DepartmentsRepository
}

// NewTicketHandler creates a new TicketHandler instance
func NewTicketHandler(
	ticketRepo *repository.TicketRepository,
	ticketTypeRepo *repository.TicketTypesRepository,
	ticketPriorityRepo *repository.TicketPrioritiesRepository,
	ticketStatusRepo *repository.TicketStatusesRepository,
	userRepo *repository.UsersRepository,
	departmentRepo *repository.DepartmentsRepository,
) *TicketHandler {
	return &TicketHandler{
		TicketRepo:         ticketRepo,
		TicketTypeRepo:     ticketTypeRepo,
		TicketPriorityRepo: ticketPriorityRepo,
		TicketStatusRepo:   ticketStatusRepo,
		UserRepo:           userRepo,
		DepartmentRepo:     departmentRepo,
	}
}

// CreateTicketHandler handles POST /tickets/CreateTicket/
// @Summary Create a new ticket
// @Description Creates a new ticket with the provided data
// @Tags Ticket
// @Accept json
// @Produce json
// @Param ticket body dto.TicketCreateRequest true "Ticket data"
// @Success 201 {object} dto.IDResponse[string]
// @Failure 400 {object} errx.APIError
// @Failure 409 {object} errx.APIError
// @Failure 500 {object} errx.APIError
// @Router /tickets/CreateTicket/ [post]
func (h *TicketHandler) CreateTicketHandler(c *gin.Context) {
	var ticketDTO dto.TicketCreateRequest

	// check valid request
	if err := c.ShouldBindJSON(&ticketDTO); err != nil {
		appErr := errx.Respond(errx.ErrBadRequest, err)
		c.JSON(appErr.HTTPStatus, appErr)
		return
	}

	// Check attachment limit
	if len(ticketDTO.Attachments) > 0 {
		total := len(ticketDTO.Attachments)
		if total > config.Get().TicketConfig.MaxTicketUploadFile {
			apiErr := errx.Respond(errx.ErrMaxTicketFilesExceeded, errors.New(""))
			apiErr.Err.Message += fmt.Sprintf(" حداکثر فایل مجاز برای هر تیکت: %d", config.Get().TicketConfig.MaxTicketUploadFile)
			c.JSON(apiErr.HTTPStatus, apiErr)
			return
		}
	}

	// check user exists
	isUserExist, err := h.UserRepo.IsUserExist(c.Request.Context(), int64(ticketDTO.UserID))
	if err != nil {
		c.JSON(err.HTTPStatus, err)
		return
	}
	if !isUserExist {
		err := errx.Respond(errx.ErrUserNotFound, errors.New("user not found"))
		c.JSON(err.HTTPStatus, err)
		return // Add return!
	}

	// check ticket type isTicketTypeExists
	isTicketTypeExists, err := h.TicketTypeRepo.IsTicketTypeExits(c.Request.Context(), int64(ticketDTO.TicketTypeID))
	if err != nil {
		c.JSON(err.HTTPStatus, err)
		return
	}
	if !isTicketTypeExists {
		err := errx.Respond(errx.ErrTicketTypeNotFound, errors.New("ticket type not found"))
		c.JSON(err.HTTPStatus, err)
		return // Add return!
	}

	// check department exists
	isDepExists, err := h.DepartmentRepo.IsDepartmentExits(c.Request.Context(), int64(ticketDTO.DepartmentID))
	if err != nil {
		c.JSON(err.HTTPStatus, err)
		return
	}
	if !isDepExists {
		err := errx.Respond(errx.ErrDepartmentNotFound, errors.New("ticket type not found"))
		c.JSON(err.HTTPStatus, err)
		return // Add return!
	}

	openStatus, err := h.TicketStatusRepo.GetOpenStatus(c.Request.Context())
	if err != nil {
		c.JSON(err.HTTPStatus, err)
		return
	}
	ticketDTO.TicketStatusID = openStatus.ID

	createdTicket, err := h.TicketRepo.CreateTicket(c.Request.Context(), &ticketDTO)
	if err != nil {
		c.JSON(err.HTTPStatus, err)
		return
	}

	c.JSON(http.StatusCreated, createdTicket)
}

// GetTicketByTrackCodeHandler handles POST /tickets/GetTicketByTrackCode/
// @Summary Get ticket by track code
// @Description Returns a ticket by its track code
// @Tags Ticket
// @Accept json
// @Produce json
// @Param request body dto.TicketByTrackCodeRequestDTO true "Track Code Request"
// @Success 200 {object} dto.TicketResponse
// @Failure 400 {object} errx.APIError
// @Failure 404 {object} errx.APIError
// @Failure 500 {object} errx.APIError
// @Router /tickets/GetTicketByTrackCode/ [post]
func (h *TicketHandler) GetTicketByTrackCodeHandler(c *gin.Context) {
	var req dto.TicketByTrackCodeRequestDTO

	// Bind JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := errx.Respond(errx.ErrBadRequest, err)
		c.JSON(appErr.HTTPStatus, appErr)
		return
	}

	_, parseErr := util.ParsTrackCode(req.TrackCode)
	if parseErr != nil {
		appErr := errx.Respond(errx.ErrBadRequest, parseErr)
		c.JSON(appErr.HTTPStatus, appErr)
		return
	}

	// Get user by username
	user, err := h.UserRepo.GetUserByUsername(c.Request.Context(), req.Username)
	if err != nil {
		if err.Err.Code == errx.ErrUserNotFound {
			err = errx.Respond(errx.ErrTicketNotFound, errors.New("username not found"))
		}
		c.JSON(err.HTTPStatus, err)
		return
	}

	// Get ticket by track code
	ticketDTO, err := h.TicketRepo.GetTicketByTrackCode(c.Request.Context(), req.TrackCode)
	if err != nil {
		c.JSON(err.HTTPStatus, err)
		return
	}

	// Ensure the ticket belongs to the user
	if ticketDTO.UserID != user.ID {
		appErr := errx.Respond(errx.ErrTicketNotFound, errors.New("this username did not create this ticket"))
		c.JSON(appErr.HTTPStatus, appErr)
		return
	}

	// Return the ticket
	c.JSON(http.StatusOK, ticketDTO)
}

// GetTicketByIDHandler handles POST /tickets/GetTicketByID/
// @Summary Get ticket by ID
// @Description Returns a ticket by its ID
// @Tags Ticket
// @Accept json
// @Produce json
// @Param request body dto.TicketByIDRequestDTO true "Ticket ID Request"
// @Success 200 {object} dto.TicketResponse
// @Failure 400 {object} errx.APIError
// @Failure 404 {object} errx.APIError
// @Failure 500 {object} errx.APIError
// @Router /tickets/GetTicketByID/ [post]
func (h *TicketHandler) GetTicketByIDHandler(c *gin.Context) {
	var req dto.TicketByIDRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil || req.ID == "" {
		appErr := errx.Respond(errx.ErrBadRequest, err)
		c.JSON(appErr.HTTPStatus, appErr)
		return
	}

	ticketDTO, err := h.TicketRepo.GetTicketByID(c.Request.Context(), req.ID)
	if err != nil {
		c.JSON(err.HTTPStatus, err)
		return
	}

	c.JSON(http.StatusOK, ticketDTO)
}

// GetTicketsListHandler handles POST /tickets/GetTicketsList/
// @Summary List tickets with paging and filtering
// @Description Returns a paginated list of tickets based on complex filter and sort options
// @Tags Ticket
// @Accept json
// @Produce json
// @Param request body dto.TicketQueryParams true "Ticket filter and paging options"
// @Success 200 {object} dto.PagingResponse[dto.TicketResponse]
// @Failure 400 {object} errx.APIError
// @Failure 500 {object} errx.APIError
// @Router /tickets/GetTicketsList/ [post]
func (h *TicketHandler) GetTicketsListHandler(c *gin.Context) {
	var req dto.TicketQueryParams

	// Bind JSON body for POST
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := errx.Respond(errx.ErrBadRequest, err)
		c.JSON(appErr.HTTPStatus, appErr)
		return
	}

	// Fetch tickets from repository
	ticketsListDTO, err := h.TicketRepo.GetTickets(c.Request.Context(), req)
	if err != nil {
		c.JSON(err.HTTPStatus, err)
		return
	}

	c.JSON(http.StatusOK, ticketsListDTO)
}

// GetAllActiveTicketTypesHandler handles GET /tickets/GetAllActiveTicketTypes/
// @Summary Get all active ticket types
// @Description Returns a list of all active ticket types
// @Tags Ticket
// @Accept json
// @Produce json
// @Success 200 {object} dto.TicketTypeDto
// @Failure 500 {object} errx.APIError
// @Router /tickets/GetAllActiveTicketTypes/ [get]
func (h *TicketHandler) GetAllActiveTicketTypesHandler(c *gin.Context) {

	ticketTypesList, err := h.TicketTypeRepo.GetAllActiveTicketTypes(c.Request.Context())
	if err != nil {
		c.JSON(err.HTTPStatus, err)
		return
	}

	ticketTypesDTO := make([]*dto.TicketTypeDto, len(ticketTypesList)-1)
	for _, v := range ticketTypesList {
		println(ticketTypesDTO)
		ticketTypesDTO = append(ticketTypesDTO, dto.ToTicketTypeDTO(&v))
	}

	c.JSON(http.StatusOK, ticketTypesDTO)
}

// GetAllActiveTicketStatusesHandler handles GET /tickets/GetAllActiveTicketStatuses/
// @Summary Get all active ticket statuses
// @Description Returns a list of all active ticket statuses
// @Tags Ticket
// @Accept json
// @Produce json
// @Success 200 {object} dto.TicketStatusDTO
// @Failure 500 {object} errx.APIError
// @Router /tickets/GetAllActiveTicketStatuses/ [get]
func (h *TicketHandler) GetAllActiveTicketStatusesHandler(c *gin.Context) {
	var ticketStatusDTO []dto.TicketStatusDTO
	ticketStatuses, err := h.TicketStatusRepo.GetAllActiveTicketStatuses(c.Request.Context())
	if err != nil {
		c.JSON(err.HTTPStatus, err)
		return
	}

	for _, v := range ticketStatuses {
		ticketStatusDTO = append(ticketStatusDTO, *dto.ToTicketStatusDTO(&v))
	}

	c.JSON(http.StatusOK, ticketStatusDTO)
}

// CloseTicketHandler handles POST /tickets/CloseTicket
// @Summary POST CloseTicket By ID
// @Description Returns a list of all active ticket statuses
// @Tags Ticket
// @Accept json
// @Produce json
// @Success 200 {object} dto.TicketStatusDTO
// @Failure 500 {object} errx.APIError
// @Router /tickets/GetAllActiveTicketStatuses/ [get]
func (h *TicketHandler) CloseTicketHandler(c *gin.Context) {
	var req dto.IDRequest[string]
	if !bindJSON(c, &req) {
		return
	}

	close, err := h.TicketStatusRepo.GetCloseStatus(c.Request.Context())
	if err != nil {
		c.JSON(err.HTTPStatus, err)
		return
	}
	ticket, err := h.TicketRepo.SetTicketStatus(c.Request.Context(), req.ID, close.ID)

	if err != nil {
		c.JSON(err.HTTPStatus, err)
		return
	}

	c.JSON(http.StatusOK, ticket)
}
