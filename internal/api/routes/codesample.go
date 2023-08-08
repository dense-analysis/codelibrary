package routes

import (
	"time"

	"github.com/dense-analysis/codelibrary/internal/api/apisession"
	"github.com/dense-analysis/codelibrary/internal/api/database"
	"github.com/dense-analysis/codelibrary/internal/api/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func validateCodeSampleSearch(c *fiber.Ctx, search *models.CodeSampleSearch) (error, bool) {
	if err := c.QueryParser(search); err != nil {
		return err, false
	}

	queries := c.Queries()

	if _, ok := queries["page"]; !ok {
		search.Page = 1
	}

	if _, ok := queries["pageSize"]; !ok {
		search.PageSize = 20
	}

	errorDetail := []models.ErrorLocation{}

	if search.Page == 0 {
		errorDetail = append(
			errorDetail,
			models.NewErrorLocation("invalidValue", "Invalid page", "query", "page"),
		)
	}

	if search.PageSize < 1 || search.PageSize > 50 {
		errorDetail = append(
			errorDetail,
			models.NewErrorLocation("invalidValue", "Invalid pageSize", "query", "pageSize"),
		)
	}

	if len(errorDetail) > 0 {
		err := sendError(c, 422, errorDetail)

		return err, false
	}

	return nil, true
}

// ListCodeSamplesHandler godoc
// @Tags Code Samples
// @Summary List Code Samples
// @Description Retrieve a list of Code Samples
// @Param q query string false "A string for searching for code samples"
// @Param l query string false "Search for results for a particular language by name"
// @Param page query integer false "The page to list results from"
// @Param pageSize query integer false "The amount of items to fetch in a given page"
// @Success 200 {object} CodeSamplePage
// @Router /code [get]
func ListCodeSamplesHandler(db database.DatabaseAPI) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var search models.CodeSampleSearch

		if err, ok := validateCodeSampleSearch(c, &search); err != nil || !ok {
			return err
		}

		page, err := db.FindCodeSamples(c.Context(), search)

		if err != nil {
			return err
		}

		return c.JSON(page)
	}
}

type CodeSampleParams struct {
	ID string `json:"id"`
}

func submitCodeSample(db database.DatabaseAPI, c *fiber.Ctx, mode SubmitMode) error {
	var err error
	var submission models.CodeSampleSubmission

	if err := c.BodyParser(&submission); err != nil {
		return err
	}

	var id uuid.UUID

	if mode == Create {
		id, err = uuid.NewRandom()
	} else {
		var params CodeSampleParams
		err = c.ParamsParser(&params)
		id, err = uuid.Parse(params.ID)

		if err != nil {
			return sendError(c, 400, []models.ErrorLocation{
				models.NewErrorLocation("invalidId", "invalid UUID", "params", "id"),
			})
		}
	}

	if err != nil {
		return err
	}

	user, err := apisession.LoadUser(c, db)

	if err != nil {
		return err
	}

	language, err := db.GetLanguage(c.Context(), submission.LanguageID)

	if err != nil {
		return err
	}

	var sample models.CodeSample

	if mode == Create {
		sample.ID = id
		sample.SubmittedBy = user
		sample.Created = time.Now()
	} else {
		sample, err = db.GetCodeSample(c.Context(), id)

		if err != nil {
			return err
		}

		if sample.SubmittedBy.ID != user.ID {
			return sendBodyError(c, 403, "forbidden", "Not your code sample")
		}
	}

	sample.Language = language
	sample.Title = submission.Title
	sample.Description = submission.Description
	sample.Body = submission.Body
	sample.Modified = time.Now()

	if mode == Create {
		err = db.CreateCodeSample(c.Context(), sample)
	} else {
		err = db.UpdateCodeSample(c.Context(), sample)
	}

	if err != nil {
		return err
	}

	if mode == Create {
		c.Status(201)
	}

	return c.JSON(sample)
}

// CreateCodeSampleHandler godoc
// @Tags Code Samples
// @Summary Submit Code Sample
// @Description Submit a new Code Sample
// @Param data body CodeSampleSubmission true "CodeSample data"
// @Success 201 {object} CodeSample
// @Router /code [post]
func CreateCodeSampleHandler(db database.DatabaseAPI) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return submitCodeSample(db, c, Create)
	}
}

// UpdateCodeSampleHandler godoc
// @Tags Code Samples
// @Summary Update a Code Sample
// @Description Update an existing Code Sample
// @Param id path string true "The UUID of the code sample to update"
// @Param data body CodeSampleSubmission true "CodeSample data"
// @Success 200 {object} CodeSample
// @Router /code/{id} [put]
func UpdateCodeSampleHandler(db database.DatabaseAPI) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return submitCodeSample(db, c, Update)
	}
}

// GetCodeSampleHandler godoc
// @Tags Code Samples
// @Summary Get a Code Sample
// @Description Get a Code Sample
// @Param id path string true "The UUID of the code sample to get"
// @Success 200 {object} CodeSample
// @Router /code/{id} [get]
func GetCodeSampleHandler(db database.DatabaseAPI) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := parseParamsID(c)

		if err != nil {
			return err
		}

		sample, err := db.GetCodeSample(c.Context(), id)

		if err != nil {
			return err
		}

		return c.JSON(sample)
	}
}

// DeleteCodeSampleHandler godoc
// @Tags Code Samples
// @Summary Delete a Code Sample
// @Description Delete a Code Sample
// @Param id path string true "The UUID of the code sample to delete"
// @Success 204
// @Router /code/{id} [delete]
func DeleteCodeSampleHandler(db database.DatabaseAPI) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := parseParamsID(c)

		if err != nil {
			return sendError(c, 400, []models.ErrorLocation{
				models.NewErrorLocation("invalidId", "invalid UUID", "params", "id"),
			})
		}

		user, err := apisession.LoadUser(c, db)

		if err != nil {
			return err
		}

		sample, err := db.GetCodeSample(c.Context(), id)

		if err != nil {
			return err
		}

		if sample.SubmittedBy.ID != user.ID {
			return sendError(c, 403, []models.ErrorLocation{
				models.NewErrorLocation("forbidden", "Not your code sample", "params", "id"),
			})
		}

		c.Status(204)
		return db.DeleteCodeSample(c.Context(), id)
	}
}
