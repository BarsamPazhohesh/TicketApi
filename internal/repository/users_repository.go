package repository

import (
	"context"
	"ticket-api/internal/db/users"
)

type UsersRepository struct {
	queries *users.Queries
}

func NewUsersRepository(queries *users.Queries) *UsersRepository {
	return &UsersRepository{
		queries: queries,
	}
}

func (repo *UsersRepository) IsUserExist(ctx context.Context, userId int64) (bool, error) {
	_, err := repo.queries.CheckUserByID(ctx, userId)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (repo *UsersRepository) AddUser(ctx context.Context, param users.CreateUserParams) (int64, error) {
	userId, err := repo.queries.CreateUser(ctx, param)
	if err != nil {
		return -1, err
	}

	return userId, nil
}

func (repo *UsersRepository) GetUserByUsername(ctx context.Context, username string) (int64, error) {
	return repo.queries.GetUserByUsername(ctx, username)
}
