package repository

import (
	"context"
	"ticket-api/internal/db/version"
	"ticket-api/internal/dto"
)

// VersionRepository handles version-related queries
type VersionRepository struct {
	queries *version.Queries
}

// NewVersionRepository creates a new VersionRepository
func NewVersionRepository(queries *version.Queries) *VersionRepository {
	return &VersionRepository{
		queries: queries,
	}
}

// GetCurrentVersion returns the current version for a given API version
func (r *VersionRepository) GetCurrentVersion(ctx context.Context, apiVersion string) (version.AppVersion, error) {
	v, err := r.queries.GetCurrentVersion(ctx, apiVersion)
	if err != nil {
		return version.AppVersion{}, err
	}
	return v, nil
}

// CreateVersion inserts a new version into app_versions using VersionDTO
func (r *VersionRepository) CreateVersion(ctx context.Context, dto dto.VersionDTO) (version.AppVersion, error) {
	model := dto.ToModel()

	params := version.CreateVersionParams{
		Version:    model.Version,
		ApiVersion: model.ApiVersion,
		Notes:      model.Notes,
		IsCurrent:  model.IsCurrent,
	}

	return r.queries.CreateVersion(ctx, params)
}

// ListVersions returns all versions for a given API version
func (r *VersionRepository) ListVersions(ctx context.Context, apiVersion string) ([]version.AppVersion, error) {
	return r.queries.ListVersions(ctx, apiVersion)
}
