package routes_test

import (
	"testing"
	"time"

	"github.com/dense-analysis/codelibrary/internal/api/apisession"
	"github.com/dense-analysis/codelibrary/internal/api/database"
	"github.com/dense-analysis/codelibrary/internal/api/models"
	"github.com/dense-analysis/codelibrary/internal/api/routes"
	"github.com/dense-analysis/codelibrary/internal/testutils"
	"github.com/dense-analysis/ranges"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestListCodeSamples(t *testing.T) {
	t.Parallel()

	r := NewRouteTester(t)
	defer r.Release()

	expectedPage := models.CodeSamplePage{
		Count:   1,
		Results: []models.CodeSample{{ID: testutils.UUIDFromInt(1)}},
	}
	r.DB.FindCodeSamplesResult.A = expectedPage

	r.AssertStatus(routes.ListCodeSamplesHandler, 200)

	var actualPage models.CodeSamplePage
	r.GetResponse(&actualPage)
	assert.Equal(t, expectedPage, actualPage)
}

func TestListCodeSamplesInvalidParams(t *testing.T) {
	var tests = map[string]struct {
		search        models.CodeSampleSearch
		expectedError models.ErrorLocation
	}{
		"InvalidSmallPageSize": {
			search:        models.CodeSampleSearch{Page: 1, PageSize: 0},
			expectedError: models.NewErrorLocation("invalidValue", "Invalid pageSize", "query", "pageSize"),
		},
		"InvalidLargePageSize": {
			search:        models.CodeSampleSearch{Page: 1, PageSize: 51},
			expectedError: models.NewErrorLocation("invalidValue", "Invalid pageSize", "query", "pageSize"),
		},
		"InvalidPage": {
			search:        models.CodeSampleSearch{Page: 0, PageSize: 20},
			expectedError: models.NewErrorLocation("invalidValue", "Invalid page", "query", "page"),
		},
	}

	for name, testData := range tests {
		testData := testData
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			r := NewRouteTester(t)
			defer r.Release()

			r.SetQueryArgs(testData.search)
			r.AssertStatus(routes.ListCodeSamplesHandler, 422)

			r.AssertResponseError(testData.expectedError)
		})
	}
}

func TestSubmitCodeSample(t *testing.T) {
	var tests = map[string]struct {
		mode routes.SubmitMode
	}{
		"POST": {mode: routes.Create},
		"PUT":  {mode: routes.Update},
	}

	for name, testData := range tests {
		testData := testData
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			r := NewRouteTester(t)
			defer r.Release()

			user := models.User{ID: testutils.UUIDFromInt(1)}
			language := models.Language{
				ID:   "python",
				Name: "Python",
			}

			// Mock session and database results.
			apisession.SaveUser(r.Ctx, user)
			r.DB.GetUserResult.A = user
			r.DB.GetLanguageResult.A = language

			submission := models.CodeSampleSubmission{
				LanguageID:  language.ID,
				Title:       "My Code Sample",
				Description: "My Description",
				Body:        "2 + 2 == 4",
			}
			r.SetRequestBody(submission)

			expectedSample := models.CodeSample{
				ID:          testutils.UUIDFromInt(0),
				SubmittedBy: user,
				Language:    language,
				Title:       submission.Title,
				Description: submission.Description,
				Body:        submission.Body,
			}

			pastTime := time.Now().Add(-1 * time.Second)

			if testData.mode == routes.Create {
				r.AssertStatus(routes.CreateCodeSampleHandler, 201)
			} else {
				expectedSample.Created = pastTime
				r.DB.GetCodeSampleResult.A = expectedSample
				r.SetParams(
					ranges.MakePair("id", testutils.UUIDFromInt(1).String()),
				)
				r.AssertStatus(routes.UpdateCodeSampleHandler, 200)
			}

			var actualSample models.CodeSample
			r.GetResponse(&actualSample)

			// Check that times are kept or shifted correctly.
			if testData.mode == routes.Create {
				assert.Greater(t, actualSample.Created.Unix(), pastTime.Unix())
			} else {
				assert.Equal(t, actualSample.Created.Unix(), pastTime.Unix())
			}

			assert.Greater(t, actualSample.Modified.Unix(), pastTime.Unix())

			// Zero out times as we compare them separately.
			expectedSample.Created = time.Time{}
			expectedSample.Modified = time.Time{}
			actualSample.Created = time.Time{}
			actualSample.Modified = time.Time{}
			actualSample.ID = testutils.UUIDFromInt(0)
			assert.Equal(t, expectedSample, actualSample)

			// Check if the data is saved to the DB is correct
			var calls [][]any

			if testData.mode == routes.Create {
				calls = r.DB.GetCalls("CreateCodeSample")
			} else {
				calls = r.DB.GetCalls("UpdateCodeSample")
			}

			assert.Equal(t, 1, len(calls))

			if len(calls) == 1 {
				actualDBSample := calls[0][0].(models.CodeSample)
				actualDBSample.ID = testutils.UUIDFromInt(0)
				actualDBSample.Created = time.Time{}
				actualDBSample.Modified = time.Time{}
				assert.Equal(t, expectedSample, actualDBSample)
			}
		})
	}
}

func TestUpdateCodeSampleValidation(t *testing.T) {
	var tests = map[string]struct {
		submission         models.CodeSampleSubmission
		paramsID           string
		sample             models.CodeSample
		expectedStatusCode int
		expectedError      models.ErrorLocation
	}{
		"WrongUser": {
			submission: models.CodeSampleSubmission{},
			paramsID:   testutils.UUIDFromInt(1).String(),
			sample: models.CodeSample{
				SubmittedBy: models.User{ID: testutils.UUIDFromInt(1)},
			},
			expectedStatusCode: 403,
			expectedError:      models.NewErrorLocation("forbidden", "Not your code sample", "body"),
		},
		"InvalidUUID": {
			submission:         models.CodeSampleSubmission{},
			paramsID:           "x",
			sample:             models.CodeSample{},
			expectedStatusCode: 400,
			expectedError:      models.NewErrorLocation("invalidId", "invalid UUID", "params", "id"),
		},
	}

	for name, testData := range tests {
		testData := testData
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			r := NewRouteTester(t)
			defer r.Release()

			user := models.User{ID: testutils.UUIDFromInt(2)}
			apisession.SaveUser(r.Ctx, user)
			r.DB.GetUserResult.A = user

			r.SetRequestBody(testData.submission)

			r.DB.GetCodeSampleResult.A = testData.sample
			r.SetParams(ranges.MakePair("id", testData.paramsID))

			r.AssertStatus(routes.UpdateCodeSampleHandler, testData.expectedStatusCode)
			r.AssertResponseError(testData.expectedError)
		})
	}
}

func TestGetCodeSample(t *testing.T) {
	t.Parallel()

	r := NewRouteTester(t)
	defer r.Release()

	expectedSample := models.CodeSample{
		ID: testutils.UUIDFromInt(1),
	}

	r.DB.GetCodeSampleResult.A = expectedSample
	r.SetParams(ranges.MakePair("id", expectedSample.ID.String()))

	r.AssertStatus(routes.GetCodeSampleHandler, 200)

	var actualSample models.CodeSample
	r.GetResponse(&actualSample)
	assert.Equal(t, expectedSample, actualSample)
}

func TestGetCodeSample404(t *testing.T) {
	t.Parallel()

	r := NewRouteTester(t)
	defer r.Release()

	id := testutils.UUIDFromInt(1)

	r.DB.GetCodeSampleResult.B = database.NotFoundErr
	r.SetParams(ranges.MakePair("id", id.String()))

	r.AssertStatus(routes.GetCodeSampleHandler, 404)
}

func TestDeleteCodeSample(t *testing.T) {
	t.Parallel()

	r := NewRouteTester(t)
	defer r.Release()

	user := models.User{ID: testutils.UUIDFromInt(1)}
	err := apisession.SaveUser(r.Ctx, user)
	r.DB.GetUserResult.A = user

	assert.Nil(t, err)

	sample := models.CodeSample{
		ID:          testutils.UUIDFromInt(123),
		SubmittedBy: user,
	}
	r.DB.GetCodeSampleResult.A = sample

	r.SetParams(ranges.MakePair("id", sample.ID.String()))
	r.AssertStatus(routes.DeleteCodeSampleHandler, 204)

	calls := r.DB.GetCalls("DeleteCodeSample")
	assert.Equal(t, 1, len(calls))

	if len(calls) == 1 {
		actualID := calls[0][0].(uuid.UUID)
		assert.Equal(t, sample.ID, actualID)
	}
}

func TestDeleteCodeSampleValidation(t *testing.T) {
	var tests = map[string]struct {
		paramsID           string
		sample             models.CodeSample
		expectedStatusCode int
		expectedError      models.ErrorLocation
	}{
		"WrongUser": {
			paramsID: testutils.UUIDFromInt(1).String(),
			sample: models.CodeSample{
				SubmittedBy: models.User{ID: testutils.UUIDFromInt(1)},
			},
			expectedStatusCode: 403,
			expectedError:      models.NewErrorLocation("forbidden", "Not your code sample", "params", "id"),
		},
		"InvalidUUID": {
			paramsID:           "x",
			sample:             models.CodeSample{},
			expectedStatusCode: 400,
			expectedError:      models.NewErrorLocation("invalidId", "invalid UUID", "params", "id"),
		},
	}

	for name, testData := range tests {
		testData := testData
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			r := NewRouteTester(t)
			defer r.Release()

			user := models.User{ID: testutils.UUIDFromInt(2)}
			apisession.SaveUser(r.Ctx, user)
			r.DB.GetUserResult.A = user

			r.DB.GetCodeSampleResult.A = testData.sample
			r.SetParams(ranges.MakePair("id", testData.paramsID))

			r.AssertStatus(routes.DeleteCodeSampleHandler, testData.expectedStatusCode)
			r.AssertResponseError(testData.expectedError)
		})
	}
}

func TestDeleteCodeSample404(t *testing.T) {
	t.Parallel()

	r := NewRouteTester(t)
	defer r.Release()

	user := models.User{ID: testutils.UUIDFromInt(1)}
	apisession.SaveUser(r.Ctx, user)
	r.DB.GetUserResult.A = user

	id := testutils.UUIDFromInt(1)

	r.DB.GetCodeSampleResult.B = database.NotFoundErr
	r.SetParams(ranges.MakePair("id", id.String()))

	r.AssertStatus(routes.DeleteCodeSampleHandler, 404)
}
