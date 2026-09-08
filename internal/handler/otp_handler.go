package handler

import (
	"net/http"
	"ticket-api/internal/dto"
	_ "ticket-api/internal/errx"
	"ticket-api/internal/services/cookie"
	"ticket-api/internal/services/otp"
	"ticket-api/internal/services/token"

	"github.com/gin-gonic/gin"
)

type OTPHandler struct {
	OTPService   *otp.OTPService
	TokenService *token.TokenService
}

func NewOTPHandler(otpService *otp.OTPService, tokenService *token.TokenService) *OTPHandler {
	return &OTPHandler{
		OTPService:   otpService,
		TokenService: tokenService,
	}
}

// SendOTP handles POST /otp/send
// @Summary Send OTP verification code
// @Description Generates a verification code, caches in Redis, logs to SMS warehouse, and sends SMS via provider asynchronously
// @Tags OTP
// @Accept json
// @Produce json
// @Param payload body dto.SendOTPDTO true "Phone number"
// @Success 200 {object} dto.SendOTPResponseDTO "Verification code sent"
// @Failure 400 {object} errx.APIError
// @Failure 500 {object} errx.APIError
// @Router /otp/send [post]
func (h *OTPHandler) SendOTP(c *gin.Context) {
	var req dto.SendOTPDTO
	if !bindJSON(c, &req) {
		return
	}

	res, apiErr := h.OTPService.SendOTP(c.Request.Context(), req.PhoneNumber)
	if apiErr != nil {
		c.JSON(apiErr.HTTPStatus, apiErr)
		return
	}

	c.JSON(http.StatusOK, res)
}

// VerifyOTP handles POST /otp/verify
// @Summary Verify OTP code
// @Description Validates provided OTP code against cached code in Redis
// @Tags OTP
// @Accept json
// @Produce json
// @Param payload body dto.VerifyOTPDTO true "Phone number and OTP code"
// @Success 200 {object} dto.VerifyOTPResponseDTO "OTP is valid"
// @Failure 400 {object} errx.APIError
// @Failure 500 {object} errx.APIError
// @Router /otp/verify [post]
func (h *OTPHandler) VerifyOTP(c *gin.Context) {
	var req dto.VerifyOTPDTO
	if !bindJSON(c, &req) {
		return
	}

	res, apiErr := h.OTPService.VerifyOTP(c.Request.Context(), req.PhoneNumber, req.Code)
	if apiErr != nil {
		c.JSON(apiErr.HTTPStatus, apiErr)
		return
	}

	if h.TokenService != nil {
		tokenStr, tokenErr := h.TokenService.NewCaptchaToken(c.ClientIP(), req.PhoneNumber)
		if tokenErr != nil {
			c.JSON(tokenErr.HTTPStatus, tokenErr)
			return
		}
		cookieService := cookie.NewCaptchaCookieService()
		cookieService.Set(c, tokenStr)
	}

	c.JSON(http.StatusOK, res)
}
