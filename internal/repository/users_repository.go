package repository

import (
	"context"
	"database/sql"
	"errors"
	"ticket-api/internal/db/users"
	"ticket-api/internal/dto"
	"ticket-api/internal/errx"
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
func (repo *UsersRepository) IsUserExist(ctx context.Context, userId int64) (bool, *errx.APIError) {
	count, err := repo.queries.CheckUserByID(ctx, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, errx.Respond(errx.ErrUserNotFound, err)
		}
		return false, errx.Respond(errx.ErrInternalServerError, err)
	}
	return count != 0, nil
}

// Add a new user
func (repo *UsersRepository) AddUser(ctx context.Context, param users.CreateUserParams) (*dto.IDResponse[int64], *errx.APIError) {
	userId, err := repo.queries.CreateUser(ctx, param)
	if err != nil {
		return nil, errx.Respond(errx.ErrInternalServerError, err)
	}
	return &dto.IDResponse[int64]{ID: userId}, nil
}

// Get user by username
func (repo *UsersRepository) GetUserByUsername(ctx context.Context, username string) (*dto.UserDTO, *errx.APIError) {
	user, err := repo.queries.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errx.Respond(errx.ErrUserNotFound, err)
		}
		return nil, errx.Respond(errx.ErrInternalServerError, err)
	}

	return dto.ToUserDTO(user), nil
}

func (repo *UsersRepository) GetUserByID(ctx context.Context, userId int) (*dto.UserDTO, *errx.APIError) {
	user, err := repo.queries.GetUserByID(ctx, int64(userId))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errx.Respond(errx.ErrTicketNotFound, err)
		}
		return nil, errx.Respond(errx.ErrInternalServerError, err)
	}

	userDTO := dto.ToUserDTO(user)
	return userDTO, nil
}
