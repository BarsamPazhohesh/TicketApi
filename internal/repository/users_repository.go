package repository

import (
	"context"
	"database/sql"
	"errors"
	"os/user"
	"ticket-api/internal/apperror"
	"ticket-api/internal/db/users"
	"ticket-api/internal/dto"
)

type UsersRepository struct {
	queries *users.Queries
}

func NewUsersRepository(queries *users.Queries) *UsersRepository {
	return &UsersRepository{
		queries: queries,
	}
}

// Check if user exists by ID
func (repo *UsersRepository) IsUserExist(ctx context.Context, userId int64) (bool, *apperror.APIError) {
	_, err := repo.queries.CheckUserByID(ctx, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, apperror.Respond(apperror.ErrUserNotFound, err)
		}
		return false, apperror.Respond(apperror.ErrInternalServerError, err)
	}
	return true, nil
}

// Add a new user
func (repo *UsersRepository) AddUser(ctx context.Context, param users.CreateUserParams) (*dto.IDResponse, *apperror.APIError) {
	userId, err := repo.queries.CreateUser(ctx, param)
	if err != nil {
		return nil, apperror.Respond(apperror.ErrInternalServerError, err)
	}
	return &dto.IDResponse{ID: userId}, nil
}

// Get user by username
func (repo *UsersRepository) GetUserByUsername(ctx context.Context, username string) (*dto.UserDTO, *apperror.APIError) {
	user, err := repo.queries.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.Respond(apperror.ErrUserNotFound, err)
		}
		return nil, apperror.Respond(apperror.ErrInternalServerError, err)
	}

	return dto.ToUserDTO(user), nil
}

func (repo *UsersRepository) GetUserByID(ctx context.Context, userId int) (*dto.UserDTO, *apperror.APIError) {
	user, err := repo.queries.GetUserByID(ctx, int64(userId))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.Respond(apperror.ErrTicketNotFound, err)
		}
		return nil, apperror.Respond(apperror.ErrInternalServerError, err)
	}

	userDTO := dto.ToUserDTO(user)
	return userDTO, nil
}
