package handler

import "ticket-api/internal/repository"

type AppHandlers struct {
	Version *VersionHandler
	Ticket  *TicketHandler
}

func NewAppHandlers(repos *repository.AppRepositories) *AppHandlers {
	return &AppHandlers{
		Version: NewVersionHandler(repos.Version),
		Ticket:  NewTicketHandler(repos.Ticket),
	}
}
