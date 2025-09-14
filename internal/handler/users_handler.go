package handler

import (
	"net/http"
	"ticket-api/internal/appError"
	"ticket-api/internal/db/users"
	"ticket-api/internal/dto"
	"ticket-api/internal/repository"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	repo *repository.UsersRepository
}

func NewUserHandler(repo *repository.UsersRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

// LoginWithNoAuth handles POST /auth/loginwithnoauth
// @Summary Login or create user without authentication
// @Description If a user with the provided username and department ID exists, it returns the user's ID. Otherwise, it creates a new user and returns the new ID.
// @Tags Auth
// @Accept json
// @Produce json
// @Param login body dto.LoginWitNoAuthDTO true "Login data"
// @Success 200 {object} dto.IDResponse "User found and ID returned"
// @Success 201 {object} dto.IDResponse "New user created and ID returned"
// @Failure 400 {object} appError.Error
// @Failure 500 {object} appError.Error
// @Router /auth/LoginWithNoAuth [post]
func (handler *UserHandler) LoginWithNoAuth(c *gin.Context) {
	var loginWithNoAuthDTO dto.LoginWitNoAuthDTO

	if err := c.ShouldBindJSON(&loginWithNoAuthDTO); err != nil {
		err := appError.Respond(appError.ErrBadRequest, err)
		c.JSON(err.HTTPStatus, err)
		return
	}

	user, err := handler.repo.GetUserByUsername(c.Request.Context(), loginWithNoAuthDTO.Username)
	// no user found
	if err != nil {
		param := users.CreateUserParams{
			Username:     loginWithNoAuthDTO.Username,
			DepartmentID: loginWithNoAuthDTO.DepartmentID,
		}

		userID, err := handler.repo.AddUser(c.Request.Context(), param)
		if err != nil {
			c.JSON(err.HTTPStatus, err)
			return
		}

		c.JSON(http.StatusCreated, userID)
		return
	}

	// user found
	c.JSON(http.StatusOK, dto.IDResponse{user.ID})
}
