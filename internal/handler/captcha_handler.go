package handler

import (
	"net/http"
	"ticket-api/internal/config"
	"ticket-api/internal/dto"
	"ticket-api/internal/errx"
	"ticket-api/internal/services/captcha"
	"ticket-api/internal/services/token"

	"github.com/gin-gonic/gin"
)

type CaptchaHandler struct {
	CaptchaService *captcha.CaptchaService
	TokenService   *token.TokenService
}

// NewCaptchaHandler constructor
func NewCaptchaHandler(captchaService *captcha.CaptchaService, tokenService *token.TokenService) *CaptchaHandler {
	return &CaptchaHandler{
		CaptchaService: captchaService,
		TokenService:   tokenService,
	}
}

// GenerateCaptchaHandler handles GET /captcha/new
// @Summary Generate new captcha
// @Description Generates a new captcha and returns its ID and image (base64)
// @Tags Captcha
// @Produce json
// @Success 200 {object} dto.CaptchaResultDTO
// @Failure 500 {object} errx.Error
// @Router /captcha/new [get]
func (h *CaptchaHandler) GenerateCaptchaHandler(c *gin.Context) {
	result, err := h.CaptchaService.GenerateCaptcha()
	if err != nil {
		appErr := errx.Respond(errx.ErrInternalServerError, err)
		c.JSON(http.StatusInternalServerError, appErr)
		return
	}

	c.JSON(http.StatusOK, result)
}

// VerifyCaptchaHandler handles POST /captcha/verify
// @Summary Verify captcha
// @Description Verifies the captcha ID and user-provided answer
// @Tags Captcha
// @Accept json
// @Produce json
// @Param request body dto.CaptchaVerifyRequest true "Captcha Request"
// @Success 200 {object} dto.CaptchaVerifyRequest
// @Failure 400 {object} errx.Error
// @Failure 401 {object} errx.Error
// @Router /captcha/verify [post]
func (h *CaptchaHandler) VerifyCaptchaHandler(c *gin.Context) {
	var req dto.CaptchaVerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := errx.Respond(errx.ErrBadRequest, err)
		c.JSON(http.StatusBadRequest, appErr)
		return
	}

	if err := h.CaptchaService.VerifyCaptcha(req.ID, req.Captcha); err != nil {
		c.JSON(err.HTTPStatus, err)
		return
	}

	token, err := h.TokenService.NewCaptchaToken()
	if err != nil {
		c.JSON(err.HTTPStatus, err)
		return
	}
	c.SetCookie(
		"captcha_token", // cookie name
		token,           // cookie value
		config.Get().Captcha.ExpiredTimeToken*60, // max age in seconds
		"/",                         // path
		"",                          // domain (empty = current domain)
		config.Get().Token.Secure,   // secure (true = only send over HTTPS)
		config.Get().Token.HttpOnly, // httpOnly (cannot be accessed by JS)
	)
	c.JSON(http.StatusOK, nil)
}
