package handler

import (
	"ticket-api/internal/errx"
	"ticket-api/internal/repository"
	"ticket-api/internal/services"

	"github.com/gin-gonic/gin"
)

type AppHandlers struct {
	Version    *VersionHandler
	Ticket     *TicketHandler
	Chat       *ChatHandler
	User       *UserHandler
	Auth       *AuthHandler
	Captcha    *CaptchaHandler
	Department *DepartmentHandler
	File       *FileHandler
}

func NewAppHandlers(repos *repository.AppRepositories, services *services.AppServices) *AppHandlers {
	return &AppHandlers{
		Version:    NewVersionHandler(repos.Version),
		Ticket:     NewTicketHandler(repos.Ticket, repos.TicketTypes, repos.TicketPriorities, repos.TicketStatus, repos.Users, repos.Departments),
		Chat:       NewChatHandler(repos.Ticket, repos.ChatRepository),
		User:       NewUserHandler(repos.Users),
		Auth:       NewAuthHandler(repos.Users, services.Token),
		Captcha:    NewCaptchaHandler(services.Captcha, services.Token),
		Department: NewDepartmentHandler(repos.Departments),
		File:       NewFileHandler(services.FileStorage),
	}
}

// bindJSON is a helper to bind JSON and handle errors
func bindJSON[T any](c *gin.Context, req *T) bool {
	if err := c.ShouldBindJSON(req); err != nil {
		apiErr := errx.Respond(errx.ErrBadRequest, err)
		c.JSON(apiErr.HTTPStatus, apiErr)
		return false
	}
	return true
}
