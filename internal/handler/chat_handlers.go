package handler

import (
	"errors"
	"fmt"
	"net/http"
	"ticket-api/internal/config"
	"ticket-api/internal/dto"
	"ticket-api/internal/errx"
	"ticket-api/internal/repository"

	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	ticketRepo *repository.TicketRepository
	chatRepo   *repository.ChatRepository
}

// NewChatHandler constructor
func NewChatHandler(ticketRepo *repository.TicketRepository, chatRepo *repository.ChatRepository) *ChatHandler {
	return &ChatHandler{ticketRepo: ticketRepo, chatRepo: chatRepo}
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

	// Bind JSON body
	var chatDTO dto.ChatMessageCreateRequest
	if err := c.ShouldBindJSON(&chatDTO); err != nil {
		appErr := errx.Respond(errx.ErrBadRequest, err)
		c.JSON(appErr.HTTPStatus, appErr)
		return
	}

	// Check attachment limit
	if len(chatDTO.Attachments) > 0 {
		count, repoErr := h.ticketRepo.GetTicketAttachmentCount(c.Request.Context(), ticketID)
		if repoErr != nil {
			c.JSON(repoErr.HTTPStatus, repoErr)
			return
		}

		total := count + len(chatDTO.Attachments)
		if total > config.Get().TicketConfig.MaxTicketUploadFile {
			apiErr := errx.Respond(errx.ErrMaxTicketFilesExceeded, errors.New(""))
			apiErr.Err.Message += fmt.Sprintf(" حداکثر فایل مجاز برای هر تیکت: %d", config.Get().TicketConfig.MaxTicketUploadFile)
			c.JSON(apiErr.HTTPStatus, apiErr)
			return
		}
	}

	// Create chat message for ticket
	updatedChat, repoErr := h.chatRepo.CreateChatMessageForTicket(c.Request.Context(), ticketID, &chatDTO)
	if repoErr != nil {
		c.JSON(repoErr.HTTPStatus, repoErr)
		return
	}

	// Respond with updated chat
	c.JSON(http.StatusCreated, updatedChat)
}
