package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	_ "ticket-api/internal/dto"
	"ticket-api/internal/repository"
)

type VersionHandler struct {
	Repo *repository.VersionRepository
}

// GetCurrentVersionHandler handles GET /api/:apiVersion
// @Summary Get current API version
// @Description Returns the current version for the given API (v1, v2, etc.)
// @Tags Version
// @Param apiVersion path string true "API Version" Enums(v1,v2)
// @Produce json
// @Success 200 {object} dto.VersionDTO
// @Failure 500 {object} map[string]string
// @Router /api/{apiVersion} [get]
func (h *VersionHandler) GetCurrentVersionHandler(c *gin.Context) {
	apiVersion := c.Param("apiVersion") // captures "v1", "v2", etc.
	version, err := h.Repo.GetCurrentVersion(c.Request.Context(), apiVersion)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, version)
}
