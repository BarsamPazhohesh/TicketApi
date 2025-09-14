package dto

import "ticket-api/internal/model"

type UserDTO struct {
	ID           int64  `json:"id"`
	Username     string `json:"username"`
	DepartmentID int64  `json:"departmentId"`
	CreatedAt    string `json:"createdAt"`
	UpdatedAt    string `json:"updatedAt"`
}

// ToUserDTO maps model.User to UserDTO
func (u *UserDTO) ToModel() *model.User {
	return &model.User{
		ID:           u.ID,
		Username:     u.Username,
		DepartmentID: u.DepartmentID,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}

func ToUserDTO(u model.User) *UserDTO {
	return &UserDTO{
		ID:           u.ID,
		Username:     u.Username,
		DepartmentID: u.DepartmentID,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}
