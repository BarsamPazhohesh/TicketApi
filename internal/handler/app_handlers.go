package handler

import "ticket-api/internal/repository"

type AppHandlers struct {
	Version *VersionHandler
}

func NewAppHandlers(repos *repository.AppRepositories) *AppHandlers {
	return &AppHandlers{
		Version: &VersionHandler{Repo: repos.Version},
	}
}
