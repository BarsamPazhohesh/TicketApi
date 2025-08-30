package dto

import (
	"database/sql"
	"ticket-api/internal/db/version"
)

// VersionDTO is used for API requests/responses
type VersionDTO struct {
	ApiVersion string `json:"apiVersion"` // JSON: apiVersion
	Version    string `json:"version"`    // JSON: version
	Notes      string `json:"notes"`      // JSON: notes
	IsCurrent  bool   `json:"isCurrent"`  // JSON: isCurrent
}

// ToModel converts VersionDTO to sqlc model (AppVersion)
func (dt *VersionDTO) ToModel() *version.AppVersion {
	isCurrent := 0
	if dt.IsCurrent {
		isCurrent = 1
	}

	return &version.AppVersion{
		ApiVersion: dt.ApiVersion,
		Version:    dt.Version,
		Notes:      sql.NullString{String: dt.Notes, Valid: dt.Notes != ""},
		IsCurrent:  int64(isCurrent),
	}
}

// FromModel converts sqlc model (AppVersion) to VersionDTO
func FromModel(m *version.AppVersion) *VersionDTO {
	isCurrent := false
	if m.IsCurrent == 1 {
		isCurrent = true
	}

	notes := ""
	if m.Notes.Valid {
		notes = m.Notes.String
	}

	return &VersionDTO{
		ApiVersion: m.ApiVersion,
		Version:    m.Version,
		Notes:      notes,
		IsCurrent:  isCurrent,
	}
}
