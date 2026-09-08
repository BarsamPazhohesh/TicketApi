package handler

import (
	"net/http"
	"ticket-api/internal/dto"
	"ticket-api/internal/errx"
	"ticket-api/internal/repository"

	"github.com/gin-gonic/gin"
)

type DepartmentHandler struct {
	DepartmentRepo *repository.DepartmentsRepository
}

func NewDepartmentHandler(departmentRepo *repository.DepartmentsRepository) *DepartmentHandler {
	return &DepartmentHandler{
		DepartmentRepo: departmentRepo,
	}
}

// GetAllActiveDepartmentsHandler handles GET /departments/GetAllActiveDepartments/
// @Summary Get all active departments
// @Description Returns a list of all active departments
// @Tags Department
// @Accept json
// @Produce json
// @Success 200 {array} dto.DepartmentDTO
// @Failure 500 {object} errx.APIError
// @Router /departments/GetAllActiveDepartments/ [get]
func (h *DepartmentHandler) GetAllActiveDepartmentsHandler(c *gin.Context) {
	var departmentsListDTO []dto.DepartmentDTO
	departmentsList, err := h.DepartmentRepo.GetAllDepartments(c.Request.Context())
	if err != nil {
		appErr := errx.Respond(errx.ErrDepartmentNotFound, err)
		c.JSON(appErr.HTTPStatus, appErr)
		return
	}

	for _, v := range departmentsList {
		departmentsListDTO = append(departmentsListDTO, *dto.ToDepartmentDTO(v))
	}

	c.JSON(http.StatusOK, departmentsListDTO)
}
