package routes

import (
	"errors"

	"github.com/dense-analysis/codelibrary/internal/api/apisession"
	"github.com/dense-analysis/codelibrary/internal/api/database"
	"github.com/dense-analysis/codelibrary/internal/api/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type LoginData struct {
	Username string `json:"username"`
	Password string `json:"password" example:"password"`
} //@name LoginData

// LoginHandler godoc
// @Tags Authentication
// @Summary Log in
// @Description Log in with user credentials
// @Param data body LoginData true "Login Data"
// @Success 200 {array} User
// @Failure 403 {object} Error
// @Router /auth/login [post]
func LoginHandler(db database.DatabaseAPI) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var loginData LoginData

		if err := c.BodyParser(&loginData); err != nil {
			return err
		}

		if len(loginData.Username) == 0 || len(loginData.Password) == 0 {
			// Don't run the query if the fields are empty.
			return sendBodyError(c, 403, "invalidCredentials", "Invalid user credentials")
		}

		user, err := db.GetUserWithCredentials(c.Context(), loginData.Username, loginData.Password)

		if err != nil {
			if errors.Is(err, database.NotFoundErr) {
				return sendBodyError(c, 403, "invalidCredentials", "Invalid user credentials")
			}

			return err
		}

		err = apisession.SaveUser(c, user)

		if err != nil {
			return err
		}

		return c.JSON(user)
	}
}

// LogoutHandler godoc
// @Tags Authentication
// @Summary Log out
// @Description Clear the user from the session
// @Success 204
// @Router /auth/logout [post]
func LogoutHandler(db database.DatabaseAPI) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return apisession.DeleteUser(c)
	}
}

// RegisterHandler godoc
// @Tags Authentication
// @Summary Register a new user
// @Description Register a new user with a given password
// @Param data body RegisterUser true "User Data"
// @Success 200 {array} User
// @Failure 422 {object} Error
// @Failure 403 {object} Error
// @Router /auth/register [post]
func RegisterHandler(db database.DatabaseAPI) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var registerUser models.RegisterUser

		if err := c.BodyParser(&registerUser); err != nil {
			return err
		}

		if registerUser.Password != registerUser.ConfirmPassword {
			return sendBodyError(c, 422, "passwordMismatch", "Passwords do not match")
		}

		if len(registerUser.Password) < 8 {
			return sendBodyError(c, 422, "badPassword", "Password too short")
		}

		if len(registerUser.Password) > 64 {
			return sendBodyError(c, 422, "badPassword", "Password too long")
		}

		id, err := uuid.NewRandom()

		if err != nil {
			return err
		}

		user := models.User{
			ID:       id,
			Username: registerUser.Username,
		}
		err = db.RegisterUser(c.Context(), user, registerUser.Password)

		if errors.Is(err, database.DuplicateErr) {
			return sendBodyError(c, 403, "duplicateUser", "User already exists")
		}

		if err != nil {
			return err
		}

		return c.JSON(user)
	}
}
