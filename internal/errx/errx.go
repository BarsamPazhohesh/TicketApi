package errx

import (
	"database/sql"
	"log"
	"net/http"
	"runtime/debug"
	"sync"
	"ticket-api/internal/env"
	"time"

	"github.com/gin-gonic/gin"
)

//
// ─── ERROR CODES ──────────────────────────────────────────────────────────────
//

type ErrorCode int

const (
	ErrInternalServerError ErrorCode = iota
	ErrTicketNotFound
	ErrUnauthorized
	ErrInvalidInput
	ErrDuplicate
	ErrBadRequest
	ErrUserNotFound
	ErrTicketTypeNotFound
	ErrDepartmentNotFound
	ErrUserDuplicate
	ErrInvalidCredentials
	ErrWeakJWTSecret
	ErrIncorrectCaptcha
)

//
// ─── DATA STRUCTURES ─────────────────────────────────────────────────────────
//

// ErrorDef holds a single error definition.
type ErrorDef struct {
	Message    string
	HTTPStatus int
}

// Error represents an error payload returned to clients.
type Error struct {
	Message string    `json:"message"`
	Code    ErrorCode `json:"code"`
	Debug   string    `json:"debug,omitempty"`
	Stack   string    `json:"stack,omitempty"`
}

// APIError is a full HTTP-aware error.
type APIError struct {
	Err        *Error `json:"errors"`
	HTTPStatus int    `json:"-"`
}

// Implement Go's error interface so APIError can be used as error.
func (e *APIError) Error() string {
	if e.Err != nil {
		return e.Err.Message
	}
	return "unknown error"
}

//
// ─── REGISTRY ────────────────────────────────────────────────────────────────
//

// Registry stores error definitions and optionally fetches messages from DB.
type Registry struct {
	defs map[ErrorCode]ErrorDef
	db   *sql.DB
	mu   sync.RWMutex
}

// global read-only registry
var registry *Registry

// evaluate debug mode once at startup
var debugMode = func() bool {
	return env.GetEnvString("GIN_MODE", "debug") == "debug" || gin.Mode() == gin.DebugMode
}()

//
// ─── REGISTRY INITIALIZATION ────────────────────────────────────────────────
//

// NewRegistry creates the default registry and optionally accepts a DB connection.
func NewRegistry(db *sql.DB) *Registry {
	r := &Registry{
		defs: map[ErrorCode]ErrorDef{
			ErrInternalServerError: {"خطای داخلی سرور", http.StatusInternalServerError},
			ErrTicketNotFound:      {"تیکت پیدا نشد", http.StatusNotFound},
			ErrUnauthorized:        {"دسترسی غیرمجاز", http.StatusUnauthorized},
			ErrInvalidInput:        {"داده ورودی نامعتبر است", http.StatusBadRequest},
			ErrDuplicate:           {"رکورد تکراری است", http.StatusConflict},
			ErrBadRequest:          {"درخواست نامعتبر", http.StatusBadRequest},
			ErrUserNotFound:        {"کاربر پیدا نشد", http.StatusNotFound},
			ErrTicketTypeNotFound:  {"نوع تیکت پیدا نشد.", http.StatusNotFound},
			ErrDepartmentNotFound:  {"دپارتمان مورد نظر پیدا نشد.", http.StatusNotFound},
			ErrUserDuplicate:       {"کاربر با این مشخصات قبلاً ثبت شده است", http.StatusConflict},
			ErrInvalidCredentials:  {"نام کاربری یا رمز عبور اشتباه است", http.StatusUnauthorized},
			ErrWeakJWTSecret:       {"کلید JWT بسیار کوتاه یا ناامن است", http.StatusInternalServerError},
			ErrIncorrectCaptcha:    {"کپچا نادرست است", http.StatusUnauthorized},
		},
		db: db,
	}

	// initial load from DB
	r.loadFromDB()

	// periodic refresh
	if db != nil {
		go r.autoRefresh(5 * time.Minute)
	}

	registry = r
	return r
}

//
// ─── DB LOADING ──────────────────────────────────────────────────────────────
//

// loadFromDB fetches error messages from DB and updates cache atomically.
func (r *Registry) loadFromDB() {
	if r.db == nil {
		return
	}

	// geting error from db

	r.mu.Lock()
	// update map here
	r.mu.Unlock()
	log.Println("error messages cache refreshed from DB")
}

// autoRefresh runs loadFromDB periodically.
func (r *Registry) autoRefresh(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		r.loadFromDB()
	}
}

//
// ─── REGISTRY ACCESS ─────────────────────────────────────────────────────────
//

// Get retrieves an ErrorDef from the registry.
func (r *Registry) Get(code ErrorCode) ErrorDef {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if def, ok := r.defs[code]; ok {
		return def
	}
	return r.defs[ErrInternalServerError]
}

//
// ─── PUBLIC FUNCTIONS ────────────────────────────────────────────────────────
//

// Make creates a plain Error.
func Make(code ErrorCode, realErr error) *Error {
	def := registry.Get(code)
	errObj := &Error{
		Code:    code,
		Message: def.Message,
	}

	// add debug + stack trace only in debug mode
	if debugMode && realErr != nil {
		errObj.Debug = realErr.Error()
		errObj.Stack = string(debug.Stack())
	}

	return errObj
}

// Respond creates an APIError with HTTP status.
func Respond(code ErrorCode, realErr error) *APIError {
	def := registry.Get(code)
	return &APIError{
		Err:        Make(code, realErr),
		HTTPStatus: def.HTTPStatus,
	}
}
