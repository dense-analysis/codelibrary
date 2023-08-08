package database_test

import (
	"context"
	"testing"

	"github.com/dense-analysis/codelibrary/internal/api/database"
	"github.com/dense-analysis/codelibrary/internal/api/models"
	"github.com/dense-analysis/codelibrary/internal/testutils"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"
)

// hashMatcher accepts a password and matches a generated hash.
type hashMatcher struct {
	password string
}

func (m hashMatcher) Match(v any) bool {
	hash, ok := v.(string)

	return ok && database.CheckPasswordHash(m.password, hash)
}

func TestGetUser(t *testing.T) {
	mock, db := startDatabaseTest(t)
	defer mock.Close()

	expectedRows := pgxmock.NewRows([]string{"username"}).
		AddRow("some_user")
	mock.ExpectQuery(`SELECT .* FROM "user" WHERE id = \$1`).
		WithArgs(testutils.UUIDFromInt(1)).
		WillReturnRows(expectedRows)

	user, err := db.GetUser(context.Background(), testutils.UUIDFromInt(1))

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfilfilled expectations: %s", err)
	}

	expectedUser := models.User{
		ID:       testutils.UUIDFromInt(1),
		Username: "some_user",
	}
	assert.Nil(t, err)
	assert.Equal(t, expectedUser, user)
}

func TestGetUserWithCredentials(t *testing.T) {
	mock, db := startDatabaseTest(t)
	defer mock.Close()

	hash, _ := database.HashPassword("password123")

	expectedRows := pgxmock.NewRows([]string{"id", "password_hash"}).
		AddRow(testutils.UUIDFromInt(1), hash)
	mock.ExpectQuery(`SELECT .* FROM "user" WHERE username = \$1`).
		WithArgs("some_user").
		WillReturnRows(expectedRows)

	user, err := db.GetUserWithCredentials(
		context.Background(),
		"some_user",
		"password123",
	)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfilfilled expectations: %s", err)
	}

	expectedUser := models.User{
		ID:       testutils.UUIDFromInt(1),
		Username: "some_user",
	}
	assert.Nil(t, err)
	assert.Equal(t, expectedUser, user)
}

func TestGetUserWithCredentialsWrongPassword(t *testing.T) {
	mock, db := startDatabaseTest(t)
	defer mock.Close()

	hash, _ := database.HashPassword("password123")

	expectedRows := pgxmock.NewRows([]string{"id", "password_hash"}).
		AddRow(testutils.UUIDFromInt(1), hash)
	mock.ExpectQuery(`SELECT .* FROM "user" WHERE username = \$1`).
		WithArgs("some_user").
		WillReturnRows(expectedRows)

	_, err := db.GetUserWithCredentials(
		context.Background(),
		"some_user",
		"wrongpassword",
	)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfilfilled expectations: %s", err)
	}

	assert.Equal(t, database.NotFoundErr, err)
}

func TestRegisterUser(t *testing.T) {
	mock, db := startDatabaseTest(t)
	defer mock.Close()

	user := models.User{
		ID:       testutils.UUIDFromInt(1),
		Username: "some_user",
	}

	mock.ExpectExec(`INSERT INTO "user" \(id, username, password_hash\)`).
		WithArgs(
			testutils.UUIDFromInt(1),
			"some_user",
			hashMatcher{"password123"},
		).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	err := db.RegisterUser(
		context.Background(),
		user,
		"password123",
	)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfilfilled expectations: %s", err)
	}

	assert.Nil(t, err)
}
