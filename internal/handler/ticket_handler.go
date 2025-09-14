package handler

import (
	"errors"
	"net/http"

	"ticket-api/internal/dto"
	"ticket-api/internal/errx"
	"ticket-api/internal/repository"

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

// CreateTicketHandler handles POST tickets
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
// @Router /tickets [post]
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
} // GetTicketHandler handles GET /tickets/:id
// @Summary Get ticket by ID
// @Description Returns a ticket by its ID
// @Tags Ticket
// @Produce json
// @Param id path string true "Ticket ID"
// @Success 200 {object} dto.TicketResponse
// @Failure 404 {object} errx.APIError
// @Failure 500 {object} errx.APIError
// @Router /tickets/{id} [get]
func (h *TicketHandler) GetTicketHandler(c *gin.Context) {
	id := c.Param("id")
	ticketDTO, err := h.TicketRepo.GetTicket(c.Request.Context(), id)
	if err != nil {
		c.JSON(err.HTTPStatus, err)
		return
	}
	c.JSON(http.StatusOK, ticketDTO)
}
