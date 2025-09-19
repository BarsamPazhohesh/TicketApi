package dto

import "ticket-api/internal/model"

type UserDTO struct {
	ID           int64  `json:"id"`
	Username     string `json:"username"`
	DepartmentID int64  `json:"departmentId"`
	Password     string `json:"-"`
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
	var password string
	if u.Password.Valid {
		password = u.Password.String
	}

	return &UserDTO{
		ID:           u.ID,
		Username:     u.Username,
		DepartmentID: u.DepartmentID,
		Password:     password,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}

type LoginWitNoAuthDTO struct {
	Username     string `json:"username" binding:"required,phoneNumber"`
	DepartmentID int64  `json:"departmentId" binding:"required"`
}

type LoginWithPasswordDTO struct {
	Username string `json:"username" binding:"required,phoneNumber"`
	Password string `json:"password" binding:"required"`
}

type SigupWithPasswordDTO struct {
	Username     string `json:"username" binding:"required,phoneNumber"`
	Password     string `json:"password" binding:"required"`
	DepartmentID int64  `json:"departmentId" binding:"required"`
}

type AuthPasswordResponseDTO struct {
	UserID      int64  `json:"userId"`
	AccessToken string `json:"accessToken"`
}

type GenerateOneTimeTokenDTO struct {
	Username string `json:"username" binding:"required"`
}

type OneTimeTokenResponseDTO struct {
	Token string `json:"token"`
}
