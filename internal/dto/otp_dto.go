package dto

type SendOTPDTO struct {
	PhoneNumber string `json:"phoneNumber" binding:"required,phoneNumber"`
}

type SendOTPResponseDTO struct {
	Message string `json:"message"`
}

type VerifyOTPDTO struct {
	PhoneNumber string `json:"phoneNumber" binding:"required,phoneNumber"`
	Code        string `json:"code" binding:"required,len=6"`
}

type VerifyOTPResponseDTO struct {
	Valid   bool   `json:"valid"`
	Message string `json:"message"`
}
