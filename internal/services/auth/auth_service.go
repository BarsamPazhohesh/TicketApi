package auth

import (
	"context"
	"errors"
	"net/http"
	"ticket-api/internal/db/users"
	"ticket-api/internal/dto"
	"ticket-api/internal/errx"
	"ticket-api/internal/repository"
	"ticket-api/internal/security"
	"ticket-api/internal/services/cookie"
	"ticket-api/internal/services/token"

	"github.com/gin-gonic/gin"
)

type AuthService struct {
	userRepo           *repository.UsersRepository
	rolesRelationsRepo *repository.RolesRelationsRepository
	tokenService       *token.TokenService
	cookieService      *cookie.CookieService
}

func NewAuthService(
	userRepo *repository.UsersRepository,
	rolesRelationsRepo *repository.RolesRelationsRepository,
	tokenService *token.TokenService,
) *AuthService {
	return &AuthService{
		userRepo:           userRepo,
		rolesRelationsRepo: rolesRelationsRepo,
		tokenService:       tokenService,
		cookieService:      cookie.NewAuthCookieService(),
	}
}

// SignUpWithPassword registers user with hashed password
func (s *AuthService) SignUpWithPassword(ctx context.Context, req dto.SignUpWithPasswordDTO) (*dto.IDResponseInt64, *errx.APIError) {
	res, apiErr := s.userRepo.CreateUserWithPassword(ctx, req)
	if apiErr != nil {
		return nil, apiErr
	}
	return &dto.IDResponseInt64{ID: res.ID}, nil
}

// LoginWithPassword validates credentials and sets secure auth cookie
func (s *AuthService) LoginWithPassword(c *gin.Context, req dto.LoginWithPasswordDTO) *errx.APIError {
	ctx := c.Request.Context()

	user, apiErr := s.userRepo.GetUserByUsername(ctx, req.Username)
	if apiErr != nil {
		if apiErr.Err.Code == errx.ErrUserNotFound {
			return errx.Respond(errx.ErrInvalidCredentials, errors.New("invalid credentials"))
		}
		return apiErr
	}

	if passErr := security.CompareHashPassword(user.Password, req.Password); passErr != nil {
		return passErr
	}

	roleIDs, _ := s.rolesRelationsRepo.GetUserRoleIDs(ctx, user.ID)

	authToken, jwtErr := s.tokenService.NewAuthToken(
		token.AuthClaims{
			UserID:   user.ID,
			Username: user.Username,
			RoleIDs:  roleIDs,
		})
	if jwtErr != nil {
		return jwtErr
	}

	c.SetSameSite(http.SameSiteLaxMode)
	s.cookieService.Set(c, authToken)
	return nil
}

// LoginWithNoAuth authenticates or creates user without password
func (s *AuthService) LoginWithNoAuth(ctx context.Context, req dto.LoginWitNoAuthDTO) (*dto.IDResponseInt64, int, *errx.APIError) {
	existingUser, apiErr := s.userRepo.GetUserByUsername(ctx, req.Username)
	if apiErr != nil && apiErr.Err.Code != errx.ErrUserNotFound {
		return nil, http.StatusInternalServerError, apiErr
	}

	if existingUser != nil {
		return &dto.IDResponseInt64{ID: existingUser.ID}, http.StatusOK, nil
	}

	params := users.CreateUserParams{
		Username:     req.Username,
		DepartmentID: req.DepartmentID,
	}

	newID, apiErr := s.userRepo.AddUser(ctx, params)
	if apiErr != nil {
		return nil, http.StatusInternalServerError, apiErr
	}

	return &dto.IDResponseInt64{ID: newID.ID}, http.StatusCreated, nil
}

// GenerateSingleUseToken issues single use token for api key
func (s *AuthService) GenerateSingleUseToken(ctx context.Context, req dto.GenerateSingleUseTokenDTO) (*dto.SingleUseTokenResponseDTO, *errx.APIError) {
	_, err := s.userRepo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	tokenStr, apiErr := s.tokenService.NewOneTimeToken(req.Username)
	if apiErr != nil {
		return nil, apiErr
	}
	return &dto.SingleUseTokenResponseDTO{Token: tokenStr}, nil
}

// LoginWithOneTimeToken consumes one-time token and sets auth cookie
func (s *AuthService) LoginWithOneTimeToken(c *gin.Context, tokenString string) *errx.APIError {
	ctx := c.Request.Context()

	claims, err := s.tokenService.ParseOneTimeToken(tokenString)
	if err != nil {
		return err
	}

	user, apiErr := s.userRepo.GetUserByUsername(ctx, claims.Username)
	if apiErr != nil {
		return apiErr
	}

	roleIDs, _ := s.rolesRelationsRepo.GetUserRoleIDs(ctx, user.ID)

	jwtToken, genErr := s.tokenService.NewAuthToken(token.AuthClaims{
		UserID:   user.ID,
		Username: user.Username,
		RoleIDs:  roleIDs,
	})
	if genErr != nil {
		return genErr
	}

	c.SetSameSite(http.SameSiteLaxMode)
	s.cookieService.Set(c, jwtToken)
	return nil
}
