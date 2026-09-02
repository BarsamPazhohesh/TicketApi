package handler

import (
	"ticket-api/internal/dto"
	_ "ticket-api/internal/errx"
	"ticket-api/internal/services/auth"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	AuthService *auth.AuthService
}

// NewAuthHandler constructor
func NewAuthHandler(authService *auth.AuthService) *AuthHandler {
	return &AuthHandler{
		AuthService: authService,
	}
}

// LoginWithNoAuth handles POST /auth/LoginWithNoAuth/
// @Summary Login or create user without authentication
// @Description If a user with the provided username and department ID exists, it returns the user's ID. Otherwise, it creates a new user and returns the new ID.
// @Tags Auth
// @Accept json
// @Produce json
// @Param login body dto.LoginWitNoAuthDTO true "Login data"
// @Success 200 {object} dto.IDResponseInt64 "User found and ID returned"
// @Success 201 {object} dto.IDResponseInt64 "New user created and ID returned"
// @Failure 400 {object} errx.APIError
// @Failure 500 {object} errx.APIError
// @Router /auth/LoginWithNoAuth/ [post]
func (h *AuthHandler) LoginWithNoAuth(c *gin.Context) {
	var req dto.LoginWitNoAuthDTO
	if !bindJSON(c, &req) {
		return
	}

	res, status, apiErr := h.AuthService.LoginWithNoAuth(c.Request.Context(), req)
	if apiErr != nil {
		c.JSON(apiErr.HTTPStatus, apiErr)
		return
	}

	c.JSON(status, res)
}

// SignUpWithPassword godoc
// @Summary      Sign up with username and password
// @Description  Create a new user with username and password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        payload  body      dto.SignUpWithPasswordDTO  true  "Signup credentials"
// @Success      201      {object}  dto.IDResponseInt64
// @Failure      400      {object}  errx.APIError
// @Failure      500      {object}  errx.APIError
// @Router       /auth/SignUp/ [post]
func (h *AuthHandler) SignUpWithPassword(c *gin.Context) {
	var req dto.SignUpWithPasswordDTO
	if !bindJSON(c, &req) {
		return
	}

	res, apiErr := h.AuthService.SignUpWithPassword(c.Request.Context(), req)
	if apiErr != nil {
		c.JSON(apiErr.HTTPStatus, apiErr)
		return
	}

	c.JSON(201, res)
}

// LoginWithPassword godoc
// @Summary      Login with username and password
// @Description  Authenticate user and set secure auth cookie
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        payload  body      dto.LoginWithPasswordDTO  true  "Login credentials"
// @Success      200
// @Failure      400      {object}  errx.APIError
// @Failure      401      {object}  errx.APIError
// @Failure      500      {object}  errx.APIError
// @Router       /auth/Login/ [post]
func (h *AuthHandler) LoginWithPassword(c *gin.Context) {
	var req dto.LoginWithPasswordDTO
	if !bindJSON(c, &req) {
		return
	}

	if apiErr := h.AuthService.LoginWithPassword(c, req); apiErr != nil {
		c.JSON(apiErr.HTTPStatus, apiErr)
		return
	}

	c.JSON(200, nil)
}

// GetSingleUseToken godoc
// @Summary      Generate one-time token for a user
// @Description  Returns a one-time JWT token to authenticate on another service
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        payload  body      dto.GenerateSingleUseTokenDTO true  "Username"
// @Success      200      {object}  dto.SingleUseTokenResponseDTO
// @Failure      400      {object}  errx.APIError
// @Failure      500      {object}  errx.APIError
// @Router       /auth/GetSingleUseToken/ [post]
func (h *AuthHandler) GetSingleUseToken(c *gin.Context) {
	var req dto.GenerateSingleUseTokenDTO
	if !bindJSON(c, &req) {
		return
	}

	res, apiErr := h.AuthService.GenerateSingleUseToken(c.Request.Context(), req)
	if apiErr != nil {
		c.JSON(apiErr.HTTPStatus, apiErr)
		return
	}

	c.JSON(200, res)
}

// LoginWithOneTimeToken godoc
// @Summary      Login using a SingleUse token
// @Description  Validates a SingleUse token and returns an auth JWT
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        token  query     string  true  "SingleUse token"
// @Success      200
// @Failure      400    {object}  errx.APIError
// @Failure      401    {object}  errx.APIError
// @Failure      500    {object}  errx.APIError
// @Router       /auth/LoginWithSingleUseToken/ [get]
func (h *AuthHandler) LoginWithOneTimeToken(c *gin.Context) {
	tokenStr := c.Query("token")
	if tokenStr == "" {
		c.JSON(400, gin.H{"error": "token is required"})
		return
	}

	if apiErr := h.AuthService.LoginWithOneTimeToken(c, tokenStr); apiErr != nil {
		c.JSON(apiErr.HTTPStatus, apiErr)
		return
	}

	c.JSON(200, nil)
}
