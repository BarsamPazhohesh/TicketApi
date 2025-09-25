package handler

import (
	"errors"
	"net/http"
	"ticket-api/internal/dto"
	"ticket-api/internal/errx"
	"ticket-api/internal/repository"

	_ "ticket-api/internal/routes"

	"github.com/gin-gonic/gin"
)

// TicketHandler handles ticket-related HTTP requests
type TicketHandler struct {
	TicketRepo         *repository.TicketRepository
	TicketTypeRepo     *repository.TicketTypesRepository
	TicketPriorityRepo *repository.TicketPrioritiesRepository
	UserRepo           *repository.UsersRepository
	DepartmentRepo     *repository.DepartmentsRepository
}

// NewTicketHandler creates a new TicketHandler instance
func NewTicketHandler(
	ticketRepo *repository.TicketRepository,
	ticketTypeRepo *repository.TicketTypesRepository,
	ticketPriorityRepo *repository.TicketPrioritiesRepository,
	userRepo *repository.UsersRepository,
	departmentRepo *repository.DepartmentsRepository,
) *TicketHandler {
	return &TicketHandler{
		TicketRepo:         ticketRepo,
		TicketTypeRepo:     ticketTypeRepo,
		TicketPriorityRepo: ticketPriorityRepo,
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
	if err := c.ShouldBindJSON(&req); err != nil || req.TrackCode == "" {
		appErr := errx.Respond(errx.ErrBadRequest, err)
		c.JSON(appErr.HTTPStatus, appErr)
		return
	}

	ticketDTO, err := h.TicketRepo.GetTicketByTrackCode(c.Request.Context(), req.TrackCode)
	if err != nil {
		c.JSON(err.HTTPStatus, err)
		return
	}

	c.JSON(http.StatusOK, ticketDTO)
}

// GetTicketByIDHandler handles POST /tickets/:id/
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
// @Router /tickets/:id/ [post]
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
