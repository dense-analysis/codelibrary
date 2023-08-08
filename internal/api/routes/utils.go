package routes

import (
	"github.com/dense-analysis/codelibrary/internal/api/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func sendError(c *fiber.Ctx, statusCode int, detail []models.ErrorLocation) error {
	return c.Status(statusCode).JSON(models.NewError(detail...))
}

func sendBodyError(c *fiber.Ctx, statusCode int, _type string, msg string) error {
	return sendError(
		c,
		statusCode,
		[]models.ErrorLocation{
			models.NewErrorLocation(_type, msg, "body"),
		},
	)
}

// SubmitMode is a false for updating or creating an object.
type SubmitMode bool

const (
	Update SubmitMode = false
	Create SubmitMode = true
)

func parseParamsID(c *fiber.Ctx) (uuid.UUID, error) {
	params := struct {
		ID string `json:"id"`
	}{}
	err := c.ParamsParser(&params)

	if err != nil {
		var id uuid.UUID
		return id, err
	}

	return uuid.Parse(params.ID)
}
