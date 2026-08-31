package handler

import (
	"ticket-api/internal/dto"
	"ticket-api/internal/services/ticket"
	"ticket-api/internal/services/token"

	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	ticketService *ticket.TicketService
}

func NewChatHandler(ticketService *ticket.TicketService) *ChatHandler {
	return &ChatHandler{
		ticketService: ticketService,
	}
}

// CreateChatHandler godoc
// @Summary      Create a new chat message for a ticket
// @Description  Adds a new chat message to the specified ticket ID, uploads attachments, and returns the created message.
// @Tags         Chat
// @Accept       json
// @Produce      json
// @Param        id       path      string                      true  "Ticket ID (UUID)"
// @Param        message  body      dto.ChatMessageCreateRequest  true  "Chat message payload"
// @Success      201      {object}  dto.ChatMessageDTO
// @Failure      400      {object}  errx.APIError
// @Failure      404      {object}  errx.APIError
// @Failure      500      {object}  errx.APIError
// @Router       /tickets/{id}/CreateChat/ [post]
func (h *ChatHandler) CreateChatHandler(c *gin.Context) {
	ticketID := c.Param("id")
	if ticketID == "" {
		c.JSON(400, gin.H{"error": "ticket id is required"})
		return
	}

	var req dto.ChatMessageCreateRequest
	if !bindJSON(c, &req) {
		return
	}

	var senderID int64 = 0
	var senderType = "user"
	if val, exists := c.Get("user"); exists {
		if claims, ok := val.(*token.AuthClaims); ok {
			senderID = claims.UserID
			for _, r := range claims.RoleIDs {
				if r == 1 || r == 2 {
					senderType = "agent"
					break
				}
			}
		}
	}

	createdChat, apiErr := h.ticketService.CreateChatMessage(c.Request.Context(), ticketID, senderID, senderType, req)
	if apiErr != nil {
		c.JSON(apiErr.HTTPStatus, apiErr)
		return
	}

	c.JSON(201, createdChat)
}
