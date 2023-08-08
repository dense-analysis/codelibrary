package apisession

import (
	"errors"
	"time"

	"github.com/dense-analysis/codelibrary/internal/api/database"
	"github.com/dense-analysis/codelibrary/internal/api/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
)

var NoUserInSessionErr = errors.New("no user id in session")

func setCookie(c *fiber.Ctx, name string, value string, expires time.Time) {
	cookie := fiber.Cookie{
		Name:     name,
		Value:    value,
		HTTPOnly: true,
		Secure:   true,
		SameSite: fiber.CookieSameSiteStrictMode,
		Expires:  expires,
	}
	c.Cookie(&cookie)
}

func saveUserID(c *fiber.Ctx, userID uuid.UUID) error {
	userIDString := userID.String()
	setCookie(c, "userID", userIDString, time.Now().Add(24*time.Hour))
	// Save the ID in locals so it can be read back later on in the request.
	c.Locals("userID", userIDString)

	return nil
}

func SaveUser(c *fiber.Ctx, user models.User) error {
	return saveUserID(c, user.ID)
}

func loadUserID(c *fiber.Ctx) (uuid.UUID, error) {
	// Try to get the ID from locals first.
	userIDString, _ := c.Locals("userID").(string)

	if len(userIDString) == 0 {
		// Try to get the ID from the request cookie.
		userIDString = c.Cookies("userID")
	}

	id, err := uuid.Parse(userIDString)

	if err != nil {
		err = NoUserInSessionErr
	}

	return id, err
}

func LoadUser(c *fiber.Ctx, db database.DatabaseAPI) (models.User, error) {
	userID, err := loadUserID(c)

	if err != nil {
		var user models.User
		return user, err
	}

	return db.GetUser(c.Context(), userID)
}

func DeleteUser(c *fiber.Ctx) error {
	// A better way to clear cookies is to set them with an expiry before now.
	setCookie(c, "userID", "", fasthttp.CookieExpireDelete)
	// Clear the user ID from locals.
	c.Locals("userID", nil)

	return nil
}
