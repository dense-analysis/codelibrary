package apisession_test

import (
	"testing"

	"github.com/dense-analysis/codelibrary/internal/api/apisession"
	"github.com/dense-analysis/codelibrary/internal/api/database/databasemock"
	"github.com/dense-analysis/codelibrary/internal/api/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

func TestUserSessionWorkflow(t *testing.T) {
	t.Parallel()

	db := databasemock.New()
	db.GetUserResult.A = models.User{ID: uuid.New()}

	app := fiber.New()
	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(ctx)

	// Ensure we can't load the user at first.
	user, err := apisession.LoadUser(ctx, db)
	assert.Equal(t, apisession.NoUserInSessionErr, err)

	// Ensure we can save the user without errors.
	err = apisession.SaveUser(ctx, user)

	// Ensure we can load the user back again without errors.
	user, err = apisession.LoadUser(ctx, db)
	assert.Nil(t, err)
	assert.Equal(t, db.GetUserResult.A.ID, user.ID)

	// Ensure we can delete the user without errors.
	err = apisession.DeleteUser(ctx)
	assert.Nil(t, err)

	// Ensure we can't load the user anymore.
	user, err = apisession.LoadUser(ctx, db)
	assert.Equal(t, apisession.NoUserInSessionErr, err)
}
