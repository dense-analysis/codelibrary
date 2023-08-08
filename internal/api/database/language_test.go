package database_test

import (
	"context"
	"testing"

	"github.com/dense-analysis/codelibrary/internal/api/database"
	"github.com/dense-analysis/codelibrary/internal/api/models"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"
)

func TestGetLanguage(t *testing.T) {
	mock, db := startDatabaseTest(t)
	defer mock.Close()

	expectedRows := pgxmock.NewRows([]string{"name"}).
		AddRow("Python")
	mock.ExpectQuery(`SELECT .* FROM language WHERE id = \$1`).
		WithArgs("python").
		WillReturnRows(expectedRows)

	language, err := db.GetLanguage(context.Background(), "python")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfilfilled expectations: %s", err)
	}

	expectedLanguage := models.Language{
		ID:   "python",
		Name: "Python",
	}
	assert.Nil(t, err)
	assert.Equal(t, expectedLanguage, language)
}

func TestGetLanguageMissing(t *testing.T) {
	mock, db := startDatabaseTest(t)
	defer mock.Close()

	mock.ExpectQuery(`SELECT .* FROM language WHERE id = \$1`).
		WithArgs("python").
		WillReturnError(pgx.ErrNoRows)

	_, err := db.GetLanguage(context.Background(), "python")

	assert.Equal(t, database.NotFoundErr, err)
}
