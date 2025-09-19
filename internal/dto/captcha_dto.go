package dto

type CaptchaVerifyRequest struct {
	ID      string `json:"id"`
	Captcha string `json:"captcha"`
}

// CaptchaResult holds the captcha ID and image in base64
type CaptchaResultDTO struct {
	ID     string `json:"id"`
	Image  string `json:"image"`            // base64
	Answer string `json:"answer,omitempty"` // optional, only in debug mode
}
