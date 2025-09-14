package errx

import (
	"net/http"
	"runtime/debug"
	"ticket-api/internal/env"

	"github.com/gin-gonic/gin"
)

type Error struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Debug   string `json:"debug,omitempty"`
	Stack   string `json:"stack,omitempty"`
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
	ErrTicketTypeNotFound
	ErrDepartmentNotFound
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
	ErrTicketTypeNotFound:  {"نوع تیکت پیدا نشد.", http.StatusNotFound},
	ErrDepartmentNotFound:  {"دپارتمان مورد نظر پیدا نشد.", http.StatusNotFound},
}

// Make creates a plain Error
func Make(code int, realErr error) *Error {
	def := errors[code]
	errObj := &Error{
		Code:    code,
		Message: def.Message,
	}

	// show debug + stack trace only in debug mode
	if env.GetEnvString("GIN_MODE", "debug") == "debug" || gin.Mode() == gin.DebugMode {
		if realErr != nil {
			errObj.Debug = realErr.Error()
			errObj.Stack = string(debug.Stack())
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
