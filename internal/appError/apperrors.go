package appError

import (
	"net/http"
	"ticket-api/internal/env"

	"github.com/gin-gonic/gin"
)

type Error struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Debug   string `json:"debug,omitempty"`
}

type APIError struct {
	Error      *Error `json:"error"`
	HTTPStatus int    `json:"-"`
}

const (
	ErrInternalServerError = iota
	ErrTicketNotFound
	ErrUnauthorized
	ErrInvalidInput
	ErrDuplicate
	ErrBadRequest
	ErrUserNotFound
)

type errorDef struct {
	Message    string
	HTTPStatus int
}

// single source of truth
var errors = map[int]errorDef{
	ErrInternalServerError: {"خطای داخلی سرور", http.StatusInternalServerError},
	ErrTicketNotFound:      {"تیکت پیدا نشد", http.StatusNotFound},
	ErrUnauthorized:        {"دسترسی غیرمجاز", http.StatusUnauthorized},
	ErrInvalidInput:        {"داده ورودی نامعتبر است", http.StatusBadRequest},
	ErrDuplicate:           {"رکورد تکراری است", http.StatusConflict},
	ErrBadRequest:          {"درخواست نامعتبر", http.StatusBadRequest},
	ErrUserNotFound:        {"کاربر پیدا نشد", http.StatusNotFound},
}

// Make creates a plain Error
func Make(code int, realErr error) *Error {
	def := errors[code]
	errObj := &Error{
		Code:    code,
		Message: def.Message,
	}

	if env.GetEnvString("GIN_MODE", "debug") == "debug" || gin.Mode() == gin.DebugMode {
		if realErr != nil {
			errObj.Debug = realErr.Error()
		}
	}

	return errObj
}

// Respond creates an APIError with HTTP status
func Respond(code int, realErr error) *APIError {
	def := errors[code]
	return &APIError{
		Error:      Make(code, realErr),
		HTTPStatus: def.HTTPStatus,
	}
}

// Error implements built-in error interface
func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.Debug != "" {
		return e.Message + " | debug: " + e.Debug
	}
	return e.Message
}
