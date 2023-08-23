package errorhandler

import (
	"errors"

	"github.com/dense-analysis/codelibrary/internal/api/apisession"
	"github.com/dense-analysis/codelibrary/internal/api/database"
	"github.com/dense-analysis/codelibrary/internal/api/models"
	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	var fiberError *fiber.Error

	if errors.As(err, &fiberError) {
		code = fiberError.Code
	}

	// If we require a user in the session, return 403.
	if errors.Is(err, apisession.NoUserInSessionErr) {
		return c.Status(403).JSON(models.NewError(
			models.NewErrorLocation("permissionDenied", "Permission Denied", "body"),
		))
	}

	// If we fail to find something from the database, return 404.
	if errors.Is(err, database.NotFoundErr) {
		return c.Status(404).JSON(models.NewError(
			models.NewErrorLocation("notFound", "Not found", "path"),
		))
	}

	c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)

	return c.Status(code).SendString(err.Error())
}
