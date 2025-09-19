package handler

import (
	"errors"
	"net/http"
	"ticket-api/internal/dto"
	"ticket-api/internal/errx"
	"ticket-api/internal/repository"
	"ticket-api/internal/security"
	"ticket-api/internal/services/cookie"
	"ticket-api/internal/services/token"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	Repo         *repository.UsersRepository
	TokenService *token.TokenService
}

// NewAuthHandler constructor
func NewAuthHandler(repo *repository.UsersRepository, tokenService *token.TokenService) *AuthHandler {
	return &AuthHandler{Repo: repo, TokenService: tokenService}
}

// SigupWithPassword godoc
// @Summary      Sign up with username and password
// @Description  Create a new user with username and password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        payload  body      dto.SigupWithPasswordDTO  true  "Signup credentials"
// @Success      201      {object}  dto.IDResponse[int64]
// @Failure      400      {object}  errx.APIError
// @Failure      500      {object}  errx.APIError
// @Router       /auth/signup [post]
func (h *AuthHandler) SigupWithPassword(c *gin.Context) {
	var credential dto.SigupWithPasswordDTO
	if err := c.ShouldBindJSON(&credential); err != nil {
		appErr := errx.Respond(errx.ErrBadRequest, err)
		c.JSON(appErr.HTTPStatus, appErr)
		return
	}

	// create user
	user, err := h.Repo.CreateUserWithPassword(c.Request.Context(), credential)
	if err != nil {
		c.JSON(err.HTTPStatus, err)
		return
	}

	c.JSON(http.StatusCreated, &dto.IDResponse[int64]{ID: user.ID})
}

// LoginWithPassword godoc
// @Summary      Login with username and password
// @Description  Authenticate user and return JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        payload  body      dto.LoginWithPasswordDTO  true  "Login credentials"
// @Success      200
// @Failure      400      {object}  errx.APIError
// @Failure      401      {object}  errx.APIError
// @Failure      500      {object}  errx.APIError
// @Router       /auth/login [post]
func (h *AuthHandler) LoginWithPassword(c *gin.Context) {
	var credential dto.LoginWithPasswordDTO
	if err := c.ShouldBindJSON(&credential); err != nil {
		appErr := errx.Respond(errx.ErrBadRequest, err)
		c.JSON(appErr.HTTPStatus, appErr)
		return
	}

	// 1. Get user
	user, err := h.Repo.GetUserByUsername(c.Request.Context(), credential.Username)
	if err != nil {
		// hide whether username or password is wrong
		if err.Err.Code == errx.ErrUserNotFound {
			err = errx.Respond(errx.ErrInvalidInput, errors.New("username or password is incorrect"))
		}
		c.JSON(err.HTTPStatus, err)
		return
	}

	// 2. Compare hashed password
	if passErr := security.CompareHashPassword(user.Password, credential.Password); passErr != nil {
		c.JSON(passErr.HTTPStatus, passErr)
		return
	}

	// 3. Generate JWT token
	token, jwtErr := h.TokenService.NewAuthToken(
		token.AuthClaims{
			UserID:   user.ID,
			Username: user.Username,
			RoleIDs:  nil,
		})
	if jwtErr != nil {
		c.JSON(jwtErr.HTTPStatus, jwtErr)
		return
	}
	cookieService := cookie.NewAuthCookieService()
	cookieService.Set(c, token)
	c.JSON(http.StatusOK, nil)
}

// GenerateOneTimeToken godoc
// @Summary      Generate one-time token for a user
// @Description  Returns a one-time JWT token to authenticate on another service
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        payload  body      dto.GenerateOneTimeTokenDTO  true  "Username"
// @Success      200      {object}  dto.OneTimeTokenResponseDTO
// @Failure      400      {object}  errx.APIError
// @Failure      500      {object}  errx.APIError
// @Router       /auth/one-time-token [post]
func (h *AuthHandler) GenerateOneTimeToken(c *gin.Context) {
	var req dto.GenerateOneTimeTokenDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := errx.Respond(errx.ErrBadRequest, err)
		c.JSON(appErr.HTTPStatus, appErr)
		return
	}

	// Optional: verify user exists
	_, err := h.Repo.GetUserByUsername(c.Request.Context(), req.Username)
	if err != nil {
		c.JSON(err.HTTPStatus, err)
		return
	}

	// Generate one-time token
	token, apiErr := h.TokenService.NewOneTimeToken(req.Username)
	if apiErr != nil {
		c.JSON(apiErr.HTTPStatus, apiErr)
		return
	}

	c.JSON(http.StatusOK, &dto.OneTimeTokenResponseDTO{
		Token: token,
	})
}

// LoginWithOneTimeToken godoc
// @Summary      Login using a one-time token
// @Description  Validates a one-time token and returns an auth JWT
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        token  query     string  true  "One-time token"
// @Success      200
// @Failure      400    {object}  errx.APIError
// @Failure      401    {object}  errx.APIError
// @Failure      500    {object}  errx.APIError
// @Router       /auth/login/token [get]
func (h *AuthHandler) LoginWithOneTimeToken(c *gin.Context) {
	tokenStr := c.Query("token")
	if tokenStr == "" {
		appErr := errx.Respond(errx.ErrBadRequest, nil)
		c.JSON(appErr.HTTPStatus, appErr)
		return
	}

	// Validate one-time token
	claims, apiErr := h.TokenService.ParseOneTimeToken(tokenStr)
	if apiErr != nil {
		c.JSON(apiErr.HTTPStatus, apiErr)
		return
	}

	// Optional: verify user exists
	user, err := h.Repo.GetUserByUsername(c.Request.Context(), claims.Username)
	if err != nil {
		c.JSON(err.HTTPStatus, err)
		return
	}

	// Generate normal auth token
	authToken, jwtErr := h.TokenService.NewAuthToken(token.AuthClaims{
		UserID:   user.ID,
		Username: user.Username,
		RoleIDs:  nil,
	})
	if jwtErr != nil {
		c.JSON(jwtErr.HTTPStatus, jwtErr)
		return
	}

	cookieService := cookie.NewAuthCookieService()
	cookieService.Set(c, authToken)

	c.JSON(http.StatusOK, nil)
}
