package handler

import (
	"net/http"
	"ticket-api/internal/dto"
	"ticket-api/internal/errx"
	"ticket-api/internal/repository"

	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	Repo *repository.TicketRepository
}

// NewChatHandler constructor
func NewChatHandler(repo *repository.TicketRepository) *ChatHandler {
	return &ChatHandler{Repo: repo}
}

// CreateChatHandler handles POST /tickets/:id/CreateChat/
// @Summary Add chat message to a ticket
// @Description Adds a new chat message to an existing ticket
// @Tags Ticket
// @Accept json
// @Produce json
// @Param id path string true "Ticket ID"
// @Param chat body dto.ChatMessageCreateRequest true "Chat message data"
// @Success 201 {object} dto.ChatMessageDTO
// @Failure 400 {object} errx.Error
// @Failure 404 {object} errx.Error
// @Failure 500 {object} errx.Error
// @Router /tickets/:id/CreateChat/ [post]
func (h *ChatHandler) CreateChatHandler(c *gin.Context) {
	ticketID := c.Param("id")

	var chatDTO dto.ChatMessageCreateRequest
	if err := c.ShouldBindJSON(&chatDTO); err != nil {
		appErr := errx.Respond(errx.ErrBadRequest, err)
		c.JSON(appErr.HTTPStatus, appErr)
		return
	}

	updatedChat, err := h.Repo.CreateChatMessageForTicket(c.Request.Context(), ticketID, &chatDTO)
	if err != nil {
		c.JSON(err.HTTPStatus, err)
	}

	c.JSON(http.StatusCreated, updatedChat)
}
