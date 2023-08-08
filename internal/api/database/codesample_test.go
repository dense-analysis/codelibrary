package database_test

import (
	"context"
	"testing"
	"time"

	"github.com/dense-analysis/codelibrary/internal/api/models"
	"github.com/dense-analysis/codelibrary/internal/testutils"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"
)

var pythonCodeSample = models.CodeSample{
	ID: testutils.UUIDFromInt(1),
	SubmittedBy: models.User{
		ID:       testutils.UUIDFromInt(123),
		Username: "some_user",
	},
	Language: models.Language{
		ID:   "python",
		Name: "Python",
	},
	Title:       "Adding two numbers",
	Description: "How to add two numbers together",
	Body:        "x + y",
}

func TestFindCodeSamples(t *testing.T) {
	mock, db := startDatabaseTest(t)
	defer mock.Close()

	firstCreated := time.Now()
	firstModified := time.Now()
	secondCreated := time.Now()
	secondModified := time.Now()

	expectedCountRows := pgxmock.
		NewRows([]string{"count"}).
		AddRow(uint64(42))
	mock.ExpectQuery(`SELECT COUNT.* FROM codesample`).
		WithArgs("search phrase").
		WillReturnRows(expectedCountRows)

	expectedRows := pgxmock.
		NewRows([]string{
			"id",
			"submitted_by_id",
			"username",
			"language_id",
			"language_name",
			"title",
			"description",
			"body",
			"created",
			"modified",
		}).
		AddRow(
			testutils.UUIDFromInt(1),
			testutils.UUIDFromInt(123),
			"some_user",
			"python",
			"Python",
			"Adding two numbers",
			"How to add two numbers together",
			"x + y",
			firstCreated,
			firstModified,
		).
		AddRow(
			testutils.UUIDFromInt(2),
			testutils.UUIDFromInt(456),
			"another_user",
			"javascript",
			"JavaScript",
			"Concatenating strings",
			"How to concatenate strings together",
			"a + b",
			secondCreated,
			secondModified,
		)
	mock.ExpectQuery(`SELECT .* FROM codesample .* LIMIT`).
		WithArgs("search phrase", uint64(20), uint64(40)).
		WillReturnRows(expectedRows)

	search := models.CodeSampleSearch{
		Query:    "search phrase",
		Page:     3,
		PageSize: 20,
	}
	page, err := db.FindCodeSamples(context.Background(), search)

	assert.Nil(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfilfilled expectations: %s", err)
	}

	// Check time values separately, as we can't compare them.
	if len(page.Results) > 0 {
		assert.Equal(t, firstCreated.Unix(), page.Results[0].Created.Unix())
		assert.Equal(t, firstModified.Unix(), page.Results[0].Modified.Unix())
	}

	if len(page.Results) > 1 {
		assert.Equal(t, secondCreated.Unix(), page.Results[1].Created.Unix())
		assert.Equal(t, secondModified.Unix(), page.Results[1].Modified.Unix())
	}

	for i := 0; i < len(page.Results); i++ {
		page.Results[i].Created = time.Time{}
		page.Results[i].Modified = time.Time{}
	}

	expectedPage := models.CodeSamplePage{
		Count: 42,
		Results: []models.CodeSample{
			pythonCodeSample,
			{
				ID: testutils.UUIDFromInt(2),
				SubmittedBy: models.User{
					ID:       testutils.UUIDFromInt(456),
					Username: "another_user",
				},
				Language: models.Language{
					ID:   "javascript",
					Name: "JavaScript",
				},
				Title:       "Concatenating strings",
				Description: "How to concatenate strings together",
				Body:        "a + b",
			},
		},
	}

	assert.Equal(t, expectedPage, page)
}

func TestFindCodeSamplesLanguagesZeroResults(t *testing.T) {
	mock, db := startDatabaseTest(t)
	defer mock.Close()

	expectedCountRows := pgxmock.
		NewRows([]string{"count"}).
		AddRow(uint64(0))
	mock.ExpectQuery(`SELECT COUNT.* FROM codesample`).
		WithArgs("search phrase", []string{"python", "javascript"}).
		WillReturnRows(expectedCountRows)

	search := models.CodeSampleSearch{
		Query:     "search phrase",
		Languages: []string{"python", "javascript"},
		Page:      1,
		PageSize:  20,
	}
	page, err := db.FindCodeSamples(context.Background(), search)

	assert.Nil(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfilfilled expectations: %s", err)
	}

	expectedPage := models.CodeSamplePage{
		Count:   0,
		Results: []models.CodeSample{},
	}
	assert.Equal(t, expectedPage, page)
}

func TestGetCodeSample(t *testing.T) {
	mock, db := startDatabaseTest(t)
	defer mock.Close()

	created := time.Now()
	modified := time.Now()

	expectedRows := pgxmock.
		NewRows([]string{
			"submitted_by_id",
			"username",
			"language_id",
			"language_name",
			"title",
			"description",
			"body",
			"created",
			"modified",
		}).
		AddRow(
			testutils.UUIDFromInt(123),
			"some_user",
			"python",
			"Python",
			"Adding two numbers",
			"How to add two numbers together",
			"x + y",
			created,
			modified,
		)
	mock.ExpectQuery(`SELECT .* FROM codesample .* WHERE id = \$1`).
		WithArgs(testutils.UUIDFromInt(1)).
		WillReturnRows(expectedRows)

	sample, err := db.GetCodeSample(context.Background(), testutils.UUIDFromInt(1))
	assert.Nil(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfilfilled expectations: %s", err)
	}

	// Check time values separately, as we can't compare them.
	assert.Equal(t, created.Unix(), sample.Created.Unix())
	assert.Equal(t, modified.Unix(), sample.Modified.Unix())

	sample.Created = time.Time{}
	sample.Modified = time.Time{}
	assert.Equal(t, pythonCodeSample, sample)
}

func TestCreateCodeSample(t *testing.T) {
	mock, db := startDatabaseTest(t)
	defer mock.Close()

	mock.ExpectExec(`INSERT INTO codesample`).
		WithArgs(
			pythonCodeSample.ID,
			pythonCodeSample.SubmittedBy.ID,
			pythonCodeSample.Language.ID,
			pythonCodeSample.Title,
			pythonCodeSample.Description,
			pythonCodeSample.Body,
			pythonCodeSample.Created,
			pythonCodeSample.Modified,
		).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	err := db.CreateCodeSample(context.Background(), pythonCodeSample)
	assert.Nil(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfilfilled expectations: %s", err)
	}
}

func TestUpdateCodeSample(t *testing.T) {
	mock, db := startDatabaseTest(t)
	defer mock.Close()

	mock.ExpectExec(`UPDATE codesample`).
		WithArgs(
			pythonCodeSample.ID,
			pythonCodeSample.SubmittedBy.ID,
			pythonCodeSample.Language.ID,
			pythonCodeSample.Title,
			pythonCodeSample.Description,
			pythonCodeSample.Body,
			pythonCodeSample.Created,
			pythonCodeSample.Modified,
		).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	err := db.UpdateCodeSample(context.Background(), pythonCodeSample)
	assert.Nil(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfilfilled expectations: %s", err)
	}
}

func TestDeleteCodeSample(t *testing.T) {
	mock, db := startDatabaseTest(t)
	defer mock.Close()

	mock.ExpectExec(`DELETE FROM codesample WHERE id = \$1`).
		WithArgs(testutils.UUIDFromInt(1)).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	err := db.DeleteCodeSample(context.Background(), testutils.UUIDFromInt(1))
	assert.Nil(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfilfilled expectations: %s", err)
	}
}
