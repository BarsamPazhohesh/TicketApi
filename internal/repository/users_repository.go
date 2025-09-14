package repository

import (
	"context"
	"database/sql"
	"errors"
	appError "ticket-api/internal/appError"
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
func (repo *UsersRepository) IsUserExist(ctx context.Context, userId int64) (bool, *appError.APIError) {
	count, err := repo.queries.CheckUserByID(ctx, userId)
	if err != nil {
		return false, appError.Respond(appError.ErrInternalServerError, err)
	}
	return count != 0, nil
}

// Add a new user
func (repo *UsersRepository) AddUser(ctx context.Context, param users.CreateUserParams) (*dto.IDResponse, *appError.APIError) {
	userId, err := repo.queries.CreateUser(ctx, param)
	if err != nil {
		return nil, appError.Respond(appError.ErrInternalServerError, err)
	}
	return &dto.IDResponse{ID: userId}, nil
}

// Get user by username
func (repo *UsersRepository) GetUserByUsername(ctx context.Context, username string) (*dto.UserDTO, *appError.APIError) {
	user, err := repo.queries.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, appError.Respond(appError.ErrUserNotFound, err)
		}
		return nil, appError.Respond(appError.ErrInternalServerError, err)
	}

	return dto.ToUserDTO(user), nil
}

func (repo *UsersRepository) GetUserByID(ctx context.Context, userId int) (*dto.UserDTO, *appError.APIError) {
	user, err := repo.queries.GetUserByID(ctx, int64(userId))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, appError.Respond(appError.ErrTicketNotFound, err)
		}
		return nil, appError.Respond(appError.ErrInternalServerError, err)
	}

	userDTO := dto.ToUserDTO(user)
	return userDTO, nil
}
