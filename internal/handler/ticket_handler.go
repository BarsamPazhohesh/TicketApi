package handler

import (
	"net/http"

	"ticket-api/internal/apperror"
	"ticket-api/internal/dto"
	"ticket-api/internal/repository"

	"github.com/gin-gonic/gin"
)

// TicketHandler handles ticket-related HTTP requests
type TicketHandler struct {
	repo *repository.TicketRepository
}

// NewTicketHandler creates a new TicketHandler instance
func NewTicketHandler(repo *repository.TicketRepository) *TicketHandler {
	return &TicketHandler{repo: repo}
}

// CreateTicketHandler handles POST /tickets
// @Summary Create a new ticket
// @Description Creates a new ticket with the provided data
// @Tags Ticket
// @Accept json
// @Produce json
// @Param ticket body dto.TicketCreateRequest true "Ticket data"
// @Success 201 {object} dto.TicketIDResponse
// @Failure 400 {object} apperror.Error
// @Failure 409 {object} apperror.Error
// @Failure 500 {object} apperror.Error
// @Router /tickets [post]
func (h *TicketHandler) CreateTicketHandler(c *gin.Context) {
	var ticketDTO dto.TicketCreateRequest
	if err := c.ShouldBindJSON(&ticketDTO); err != nil {
		c.JSON(http.StatusBadRequest, apperror.Make(apperror.ErrInvalidInput, err))
		return
	}

	createdTicket, err := h.repo.CreateTicket(c.Request.Context(), &ticketDTO)
	if err != nil {
		c.JSON(err.HTTPStatus, err)
		return
	}
	c.JSON(http.StatusCreated, createdTicket)
}

// GetTicketHandler handles GET /tickets/:id
// @Summary Get ticket by ID
// @Description Returns a ticket by its ID
// @Tags Ticket
// @Produce json
// @Param id path string true "Ticket ID"
// @Success 200 {object} dto.TicketResponse
// @Failure 404 {object} apperror.Error
// @Failure 500 {object} apperror.Error
// @Router /tickets/{id} [get]
func (h *TicketHandler) GetTicketHandler(c *gin.Context) {
	id := c.Param("id")
	ticketDTO, err := h.repo.GetTicket(c.Request.Context(), id)
	if err != nil {
		c.JSON(err.HTTPStatus, err)
		return
	}
	c.JSON(http.StatusOK, ticketDTO)
}

