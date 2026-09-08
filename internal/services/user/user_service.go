package user

import (
	"context"
	"ticket-api/internal/dto"
	"ticket-api/internal/errx"
	"ticket-api/internal/repository"
)

type UserService struct {
	userRepo *repository.UsersRepository
}

func NewUserService(userRepo *repository.UsersRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) GetUserByID(ctx context.Context, userID int64) (*dto.UserDTO, *errx.APIError) {
	return s.userRepo.GetUserByID(ctx, userID)
}

func (s *UserService) GetUserByUsername(ctx context.Context, username string) (*dto.UserDTO, *errx.APIError) {
	return s.userRepo.GetUserByUsername(ctx, username)
}

func (s *UserService) GetUsersByIDs(ctx context.Context, ids []int64) ([]*dto.UserDTO, *errx.APIError) {
	return s.userRepo.GetUsersByIDs(ctx, ids)
}
