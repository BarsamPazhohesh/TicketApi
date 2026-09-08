package handler

import (
	"net/http"
	"ticket-api/internal/dto"
	_ "ticket-api/internal/errx"
	"ticket-api/internal/services/user"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *user.UserService
}

func NewUserHandler(userService *user.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// GetUsersByIDs handles POST /users/GetUsersByIDs/
// @Summary Get all users by IDs
// @Description Returns a list of all users requested
// @Tags User
// @Accept json
// @Produce json
// @Param usersIDs body dto.UserIDsDTO true "Users IDs"
// @Success 200 {array} dto.UserDTO
// @Failure 400 {object} errx.APIError
// @Failure 404 {object} errx.APIError
// @Failure 500 {object} errx.APIError
// @Router /users/GetUsersByIDs/ [post]
func (h *UserHandler) GetUsersByIDs(c *gin.Context) {
	var req dto.UserIDsDTO
	if !bindJSON(c, &req) {
		return
	}

	users, apiErr := h.userService.GetUsersByIDs(c.Request.Context(), req.IDs)
	if apiErr != nil {
		c.JSON(apiErr.HTTPStatus, apiErr)
		return
	}

	c.JSON(http.StatusOK, users)
}

// GetUserByID handles POST /users/GetUserByID/
// @Summary Get user by ID
// @Description Returns the requested user
// @Tags User
// @Accept json
// @Produce json
// @Param userID body dto.IDRequestInt64 true "User ID"
// @Success 200 {object} dto.UserDTO
// @Failure 400 {object} errx.APIError
// @Failure 404 {object} errx.APIError
// @Failure 500 {object} errx.APIError
// @Router /users/GetUserByID/ [post]
func (h *UserHandler) GetUserByID(c *gin.Context) {
	var req dto.IDRequestInt64
	if !bindJSON(c, &req) {
		return
	}

	user, apiErr := h.userService.GetUserByID(c.Request.Context(), req.ID)
	if apiErr != nil {
		c.JSON(apiErr.HTTPStatus, apiErr)
		return
	}

	c.JSON(http.StatusOK, user)
}

// GetUserByUsername handles POST /users/GetUserByUsername/
// @Summary Get user by username
// @Description Returns the requested user
// @Tags User
// @Accept json
// @Produce json
// @Param username body dto.UsernameDTO true "Username"
// @Success 200 {object} dto.UserDTO
// @Failure 400 {object} errx.APIError
// @Failure 404 {object} errx.APIError
// @Failure 500 {object} errx.APIError
// @Router /users/GetUserByUsername/ [post]
func (h *UserHandler) GetUserByUsername(c *gin.Context) {
	var req dto.UsernameDTO
	if !bindJSON(c, &req) {
		return
	}

	user, apiErr := h.userService.GetUserByUsername(c.Request.Context(), req.Username)
	if apiErr != nil {
		c.JSON(apiErr.HTTPStatus, apiErr)
		return
	}

	c.JSON(http.StatusOK, user)
}
