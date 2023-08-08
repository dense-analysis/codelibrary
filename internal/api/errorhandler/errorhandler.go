package errorhandler

import (
	"errors"

	"github.com/dense-analysis/codelibrary/internal/api/database"
	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	var fiberError *fiber.Error

	if errors.As(err, &fiberError) {
		code = fiberError.Code
	}

	if errors.Is(err, database.NotFoundErr) {
		return c.Status(404).SendString("Not found")
	}

	c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)

	return c.Status(code).SendString(err.Error())
}
