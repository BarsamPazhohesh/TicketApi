package repository

import (
	"context"
	"database/sql"
	"errors"
	"ticket-api/internal/db/users"
	"ticket-api/internal/dto"
	"ticket-api/internal/errx"

	"golang.org/x/crypto/bcrypt"
)

type UsersRepository struct {
	queries *users.Queries
}

func NewUsersRepository(queries *users.Queries) *UsersRepository {
	return &UsersRepository{
		queries: queries,
	}
}

// IsUserExist Check if user exists by ID
func (repo *UsersRepository) IsUserExist(ctx context.Context, userID int64) (bool, *errx.APIError) {
	count, err := repo.queries.CheckUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, errx.Respond(errx.ErrUserNotFound, err)
		}
		return false, errx.Respond(errx.ErrInternalServerError, err)
	}
	return count != 0, nil
}

// AddUser Add a new user
func (repo *UsersRepository) AddUser(ctx context.Context, param users.CreateUserParams) (*dto.IDResponse[int64], *errx.APIError) {
	userID, err := repo.queries.CreateUser(ctx, param)
	if err != nil {
		return nil, errx.Respond(errx.ErrInternalServerError, err)
	}
	return &dto.IDResponse[int64]{ID: userID}, nil
}

func (repo *UsersRepository) CreateUserWithPassword(ctx context.Context, credential dto.SignUpWithPasswordDTO) (*dto.IDResponse[int64], *errx.APIError) {
	// 1. Check if username already exists
	existingUser, apiErr := repo.GetUserByUsername(ctx, credential.Username)
	if apiErr != nil && apiErr.Err.Code != errx.ErrUserNotFound {
		return nil, apiErr
	}
	if existingUser != nil {
		return nil, errx.Respond(errx.ErrUserDuplicate, errors.New("username already exists"))
	}

	// 2. Hash the password before storing
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(credential.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errx.Respond(errx.ErrInternalServerError, err)
	}
	credential.Password = string(hashedPassword)

	// 3. Create the user
	params := &users.CreateUserWithPasswordParams{
		Username: credential.Username,
		Password: sql.NullString{
			String: credential.Password,
			Valid:  credential.Password != "",
		},
		DepartmentID: credential.DepartmentID,
	}
	userID, err := repo.queries.CreateUserWithPassword(ctx, *params)
	if err != nil {
		return nil, errx.Respond(errx.ErrInternalServerError, err)
	}

	return &dto.IDResponse[int64]{ID: userID}, nil
}

// GetUserByUsername Get user by username
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

func (repo *UsersRepository) GetUserByID(ctx context.Context, userID int64) (*dto.UserDTO, *errx.APIError) {
	user, err := repo.queries.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errx.Respond(errx.ErrUserNotFound, err)
		}
		return nil, errx.Respond(errx.ErrInternalServerError, err)
	}

	userDTO := dto.ToUserDTO(user)
	return userDTO, nil
}

// GetUsersByIDs Get users by users Ids
func (repo *UsersRepository) GetUsersByIDs(ctx context.Context, userIDs []int64) ([]*dto.UserDTO, *errx.APIError) {
	users, err := repo.queries.GetUsersByIDs(ctx, userIDs)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errx.Respond(errx.ErrUserNotFound, err)
		}
		return nil, errx.Respond(errx.ErrInternalServerError, err)
	}

	var usersDTO []*dto.UserDTO

	for _, user := range users {
		usersDTO = append(usersDTO, dto.ToUserDTO(user))
	}

	return usersDTO, nil
}
