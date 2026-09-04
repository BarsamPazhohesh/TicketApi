package handler

import (
	"ticket-api/internal/errx"
	"ticket-api/internal/repository"
	"ticket-api/internal/services"

	"github.com/gin-gonic/gin"
)

type AppHandlers struct {
	Ticket     *TicketHandler
	Chat       *ChatHandler
	User       *UserHandler
	Auth       *AuthHandler
	Captcha    *CaptchaHandler
	Department *DepartmentHandler
	File       *FileHandler
	OTP        *OTPHandler
}

func NewAppHandlers(repos *repository.AppRepositories, services *services.AppServices) *AppHandlers {
	return &AppHandlers{
		Ticket:     NewTicketHandler(services.Ticket),
		Chat:       NewChatHandler(services.Ticket),
		User:       NewUserHandler(services.User),
		Auth:       NewAuthHandler(services.Auth),
		Captcha:    NewCaptchaHandler(services.Captcha, services.Token),
		Department: NewDepartmentHandler(repos.Departments),
		File:       NewFileHandler(services.FileStorage),
		OTP:        NewOTPHandler(services.OTP),
	}
}

// bindJSON is a helper to bind JSON and handle errors
func bindJSON[T any](c *gin.Context, req *T) bool {
	if err := c.ShouldBindJSON(req); err != nil {
		appErr := errx.Respond(errx.ErrBadRequest, err)
		c.JSON(appErr.HTTPStatus, appErr)
		return false
	}
	return true
}
