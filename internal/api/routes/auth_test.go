package routes_test

import (
	"testing"

	"github.com/dense-analysis/codelibrary/internal/api/apisession"
	"github.com/dense-analysis/codelibrary/internal/api/database"
	"github.com/dense-analysis/codelibrary/internal/api/models"
	"github.com/dense-analysis/codelibrary/internal/api/routes"
	"github.com/dense-analysis/codelibrary/internal/testutils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	t.Parallel()

	r := NewRouteTester(t)
	defer r.Release()

	expectedUser := models.User{
		ID:       testutils.UUIDFromInt(1),
		Username: "user",
	}

	r.DB.GetUserWithCredentialsResult.A = expectedUser

	r.SetRequestBody(routes.LoginData{Username: "user", Password: "123"})
	r.AssertStatus(routes.LoginHandler, 200)

	var actualUser models.User
	r.GetResponse(&actualUser)
	assert.Equal(t, expectedUser, actualUser)

	// Ensure we save the user ID in the session correctly.
	apisession.LoadUser(r.Ctx, r.DB)
	calls := r.DB.GetCalls("GetUser")
	assert.Equal(t, 1, len(calls))

	if len(calls) == 1 {
		id := calls[0][0].(uuid.UUID)
		assert.Equal(t, expectedUser.ID, id)
	}
}

func TestLoginErrors(t *testing.T) {
	var tests = map[string]struct {
		loginData          routes.LoginData
		expectedStatusCode int
		databaseError      error
	}{
		"UserNotFound": {
			loginData:          routes.LoginData{Username: "user", Password: "123"},
			expectedStatusCode: 403,
			databaseError:      database.NotFoundErr,
		},
		"EmptyUsername": {
			loginData:          routes.LoginData{Username: "", Password: "123"},
			expectedStatusCode: 403,
		},
		"EmptyPassword": {
			loginData:          routes.LoginData{Username: "user", Password: ""},
			expectedStatusCode: 403,
		},
		"LongPassword": {
			loginData: routes.LoginData{
				Username: "user",
				Password: testutils.GenerateString('x', 65),
			},
			expectedStatusCode: 403,
		},
	}

	expectedError := models.NewErrorLocation("invalidCredentials", "Invalid user credentials", "body")

	for name, testData := range tests {
		testData := testData
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			r := NewRouteTester(t)
			defer r.Release()

			r.DB.GetUserWithCredentialsResult.B = testData.databaseError
			r.SetRequestBody(testData.loginData)

			r.AssertStatus(routes.LoginHandler, testData.expectedStatusCode)
			r.AssertResponseError(expectedError)
		})
	}
}

func TestLogout(t *testing.T) {
	t.Parallel()

	r := NewRouteTester(t)
	defer r.Release()

	user := models.User{ID: testutils.UUIDFromInt(1)}
	apisession.SaveUser(r.Ctx, user)

	r.AssertStatus(routes.LogoutHandler, 200)

	_, err := apisession.LoadUser(r.Ctx, r.DB)
	assert.Equal(t, apisession.NoUserInSessionErr, err)
}

func TestRegister(t *testing.T) {
	t.Parallel()

	r := NewRouteTester(t)
	defer r.Release()

	registerUser := models.RegisterUser{
		Username:        "user",
		Password:        "123456789",
		ConfirmPassword: "123456789",
	}
	r.SetRequestBody(registerUser)

	r.AssertStatus(routes.RegisterHandler, 200)

	expectedUser := models.User{
		Username: "user",
	}
	var actualUser models.User
	r.GetResponse(&actualUser)
	expectedUser.ID = actualUser.ID
	assert.Equal(t, expectedUser, actualUser)
}

func TestRegisterErrors(t *testing.T) {
	var tests = map[string]struct {
		data               models.RegisterUser
		expectedStatusCode int
		databaseError      error
		expectedError      models.ErrorLocation
	}{
		"PasswordMismatch": {
			data: models.RegisterUser{
				Username:        "user",
				Password:        "123456789",
				ConfirmPassword: "12345678x",
			},
			expectedStatusCode: 422,
			expectedError: models.NewErrorLocation(
				"passwordMismatch",
				"Passwords do not match",
				"body",
			),
		},
		"PasswordTooShort": {
			data: models.RegisterUser{
				Username:        "user",
				Password:        "1234567",
				ConfirmPassword: "1234567",
			},
			expectedStatusCode: 422,
			expectedError: models.NewErrorLocation(
				"badPassword",
				"Password too short",
				"body",
			),
		},
		"PasswordTooLong": {
			data: models.RegisterUser{
				Username:        "user",
				Password:        testutils.GenerateString('x', 65),
				ConfirmPassword: testutils.GenerateString('x', 65),
			},
			expectedStatusCode: 422,
			expectedError: models.NewErrorLocation(
				"badPassword",
				"Password too long",
				"body",
			),
		},
		"UsernameTooLong": {
			data: models.RegisterUser{
				Username:        testutils.GenerateString('x', 256),
				Password:        testutils.GenerateString('x', 8),
				ConfirmPassword: testutils.GenerateString('x', 8),
			},
			expectedStatusCode: 422,
			expectedError: models.NewErrorLocation(
				"badUsername",
				"Username too long",
				"body",
			),
		},
		"DuplicateUser": {
			data: models.RegisterUser{
				Username:        "user",
				Password:        "123456789",
				ConfirmPassword: "123456789",
			},
			databaseError:      database.DuplicateErr,
			expectedStatusCode: 403,
			expectedError: models.NewErrorLocation(
				"duplicateUser",
				"User already exists",
				"body",
			),
		},
	}

	for name, testData := range tests {
		testData := testData
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			r := NewRouteTester(t)
			defer r.Release()

			r.DB.RegisterUserResult = testData.databaseError

			r.SetRequestBody(testData.data)
			r.AssertStatus(routes.RegisterHandler, testData.expectedStatusCode)

			r.AssertResponseError(testData.expectedError)
		})
	}
}
