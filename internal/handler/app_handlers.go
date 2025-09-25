package handler

import (
	"ticket-api/internal/repository"
	"ticket-api/internal/services"
)

type AppHandlers struct {
	Version *VersionHandler
	Ticket  *TicketHandler
	Chat    *ChatHandler
	// User    *UserHandler
	Auth    *AuthHandler
	Captcha *CaptchaHandler
}

func NewAppHandlers(repos *repository.AppRepositories, services *services.AppServices) *AppHandlers {
	return &AppHandlers{
		Version: NewVersionHandler(repos.Version),
		Ticket:  NewTicketHandler(repos.Ticket, repos.TicketTypes, repos.TicketPriorities, repos.Users, repos.Departments),
		Chat:    NewChatHandler(repos.Ticket),
		// User:    NewUserHandler(repos.Users),
		Auth:    NewAuthHandler(repos.Users, services.Token),
		Captcha: NewCaptchaHandler(services.Captcha, services.Token),
	}
}
