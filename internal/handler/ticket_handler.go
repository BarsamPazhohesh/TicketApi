package handler

import (
	"net/http"
	"ticket-api/internal/dto"
	"ticket-api/internal/repository"

	"github.com/gin-gonic/gin"
)

type TicketHandler struct {
	Repo *repository.TicketRepository
}

// NewTicketHandler constructor
func NewTicketHandler(repo *repository.TicketRepository) *TicketHandler {
	return &TicketHandler{Repo: repo}
}

// CreateTicketHandler handles POST /tickets
// @Summary Create a new ticket
// @Description Creates a new ticket with given data
// @Tags Ticket
// @Accept json
// @Produce json
// @Param ticket body dto.TicketCreateRequest true "Ticket data"
// @Success 201 {object} dto.TicketIDResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tickets [post]
func (h *TicketHandler) CreateTicketHandler(c *gin.Context) {
	var ticketDTO dto.TicketCreateRequest
	if err := c.ShouldBindJSON(&ticketDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdTicket, err := h.Repo.CreateTicket(c.Request.Context(), &ticketDTO)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tickets/{id} [get]
func (h *TicketHandler) GetTicketHandler(c *gin.Context) {
	id := c.Param("id")
	ticketDTO, err := h.Repo.GetTicket(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if ticketDTO == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
		return
	}

	c.JSON(http.StatusOK, ticketDTO)
}

// RegisterRoutes registers ticket routes in Gin
func (h *TicketHandler) RegisterRoutes(router *gin.Engine) {
	v1 := router.Group("/api/v1")
	{
		v1.POST("/tickets", h.CreateTicketHandler)
		v1.GET("/tickets/:id", h.GetTicketHandler)
	}
}
