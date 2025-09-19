package cookie

import (
	"ticket-api/internal/config"

	"github.com/gin-gonic/gin"
)

type CookieService struct {
	Name     string
	MaxAge   int
	Path     string
	Domain   string
	Secure   bool
	HTTPOnly bool
}

// NewCookieService creates a new cookie service
func NewCookieService(name string, maxAge int) *CookieService {
	return &CookieService{
		Name:     name,
		MaxAge:   maxAge,
		Path:     "/",
		Domain:   "",
		Secure:   config.Get().Token.Secure,
		HTTPOnly: config.Get().Token.HttpOnly,
	}
}

func NewAuthCookieService() *CookieService {
	return NewCookieService(config.Get().Auth.CookieName, config.Get().Captcha.ExpiredTimeToken*60)
}

func NewCaptchaCookieService() *CookieService {
	return NewCookieService(config.Get().Captcha.CookieName, config.Get().Auth.ExpiredTimeToken*60)
}

// Set sets a cookie
func (s *CookieService) Set(c *gin.Context, value string) {
	c.SetCookie(s.Name, value, s.MaxAge, s.Path, s.Domain, s.Secure, s.HTTPOnly)
}

// Get retrieves a cookie
func (s *CookieService) Get(c *gin.Context) (string, error) {
	return c.Cookie(s.Name)
}

// Clear removes a cookie
func (s *CookieService) Clear(c *gin.Context) {
	c.SetCookie(s.Name, "", -1, s.Path, s.Domain, s.Secure, s.HTTPOnly)
}
