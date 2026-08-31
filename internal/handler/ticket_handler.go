package handler

import (
	"net/http"
	"ticket-api/internal/dto"
	"ticket-api/internal/services/ticket"
	"ticket-api/internal/services/token"

	_ "ticket-api/internal/routes"

	"github.com/gin-gonic/gin"
)

// TicketHandler handles ticket-related HTTP requests
type TicketHandler struct {
	ticketService *ticket.TicketService
}

// NewTicketHandler creates a new TicketHandler instance
func NewTicketHandler(ticketService *ticket.TicketService) *TicketHandler {
	return &TicketHandler{
		ticketService: ticketService,
	}
}

// CreateTicketHandler handles POST /tickets/CreateTicket/
// @Summary Create a new ticket
// @Description Creates a new ticket with the provided data
// @Tags Ticket
// @Accept json
// @Produce json
// @Param ticket body dto.TicketCreateRequest true "Ticket data"
// @Success 201 {object} dto.TicketCreateResponse
// @Failure 400 {object} errx.APIError
// @Failure 409 {object} errx.APIError
// @Failure 500 {object} errx.APIError
// @Router /tickets/CreateTicket/ [post]
func (h *TicketHandler) CreateTicketHandler(c *gin.Context) {
	var req dto.TicketCreateRequest
	if !bindJSON(c, &req) {
		return
	}

	var currentUserID int64 = 0
	if val, exists := c.Get("user"); exists {
		if claims, ok := val.(*token.AuthClaims); ok {
			currentUserID = claims.UserID
		}
	}

	createdTicket, apiErr := h.ticketService.CreateTicket(c.Request.Context(), currentUserID, req)
	if apiErr != nil {
		c.JSON(apiErr.HTTPStatus, apiErr)
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
	if !bindJSON(c, &req) {
		return
	}

	ticketDTO, apiErr := h.ticketService.GetTicketByTrackCode(c.Request.Context(), req.TrackCode, req.Username)
	if apiErr != nil {
		c.JSON(apiErr.HTTPStatus, apiErr)
		return
	}

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
// @Failure 403 {object} errx.APIError
// @Failure 404 {object} errx.APIError
// @Failure 500 {object} errx.APIError
// @Router /tickets/GetTicketByID/ [post]
func (h *TicketHandler) GetTicketByIDHandler(c *gin.Context) {
	var req dto.TicketByIDRequestDTO
	if !bindJSON(c, &req) {
		return
	}

	var currentUserID int64 = 0
	var roleIDs []int64
	if val, exists := c.Get("user"); exists {
		if claims, ok := val.(*token.AuthClaims); ok {
			currentUserID = claims.UserID
			roleIDs = claims.RoleIDs
		}
	}

	ticketDTO, apiErr := h.ticketService.GetTicketByID(c.Request.Context(), req.ID, currentUserID, roleIDs)
	if apiErr != nil {
		c.JSON(apiErr.HTTPStatus, apiErr)
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
// @Success 200 {object} dto.TicketPagingResponse
// @Failure 400 {object} errx.APIError
// @Failure 500 {object} errx.APIError
// @Router /tickets/GetTicketsList/ [post]
func (h *TicketHandler) GetTicketsListHandler(c *gin.Context) {
	var req dto.TicketQueryParams
	if !bindJSON(c, &req) {
		return
	}

	ticketsListDTO, apiErr := h.ticketService.GetTicketsList(c.Request.Context(), req)
	if apiErr != nil {
		c.JSON(apiErr.HTTPStatus, apiErr)
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
// @Success 200 {array} dto.TicketTypeDto
// @Failure 500 {object} errx.APIError
// @Router /tickets/GetAllActiveTicketTypes/ [get]
func (h *TicketHandler) GetAllActiveTicketTypesHandler(c *gin.Context) {
	ticketTypesDTO, apiErr := h.ticketService.GetAllActiveTicketTypes(c.Request.Context())
	if apiErr != nil {
		c.JSON(apiErr.HTTPStatus, apiErr)
		return
	}

	c.JSON(http.StatusOK, ticketTypesDTO)
}

// GetAllActiveTicketStatusesHandler handles GET /tickets/GetAllActiveTicketStatuses/
// @Summary Get all active ticket statuses
// @Description Returns a list of all active ticket statuses
// @Tags Ticket
// @Accept json
// @Produce json
// @Success 200 {array} dto.TicketStatusDTO
// @Failure 500 {object} errx.APIError
// @Router /tickets/GetAllActiveTicketStatuses/ [get]
func (h *TicketHandler) GetAllActiveTicketStatusesHandler(c *gin.Context) {
	ticketStatusDTO, apiErr := h.ticketService.GetAllActiveTicketStatuses(c.Request.Context())
	if apiErr != nil {
		c.JSON(apiErr.HTTPStatus, apiErr)
		return
	}

	c.JSON(http.StatusOK, ticketStatusDTO)
}

// CloseTicketHandler handles POST /tickets/CloseTicket/
// @Summary Close ticket by ID
// @Description Closes ticket by setting status to close
// @Tags Ticket
// @Accept json
// @Produce json
// @Param request body dto.IDRequestString true "Ticket ID"
// @Success 200 {object} dto.TicketResponse
// @Failure 400 {object} errx.APIError
// @Failure 500 {object} errx.APIError
// @Router /tickets/CloseTicket/ [post]
func (h *TicketHandler) CloseTicketHandler(c *gin.Context) {
	var req dto.IDRequestString
	if !bindJSON(c, &req) {
		return
	}

	ticket, apiErr := h.ticketService.CloseTicket(c.Request.Context(), req.ID)
	if apiErr != nil {
		c.JSON(apiErr.HTTPStatus, apiErr)
		return
	}

	c.JSON(http.StatusOK, ticket)
}
