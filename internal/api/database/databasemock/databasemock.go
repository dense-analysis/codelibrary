package databasemock

import (
	"context"

	"github.com/dense-analysis/codelibrary/internal/api/database"
	"github.com/dense-analysis/codelibrary/internal/api/models"
	"github.com/dense-analysis/ranges"
	"github.com/google/uuid"
)

type MockDatabaseAPI struct {
	calls                        map[string][][]any
	GetUserResult                ranges.Pair[models.User, error]
	GetUserWithCredentialsResult ranges.Pair[models.User, error]
	RegisterUserResult           error
	GetLanguageResult            ranges.Pair[models.Language, error]
	FindCodeSamplesResult        ranges.Pair[models.CodeSamplePage, error]
	GetCodeSampleResult          ranges.Pair[models.CodeSample, error]
	CreateCodeSampleResult       error
	UpdateCodeSampleResult       error
	DeleteCodeSampleResult       error
}

func (db *MockDatabaseAPI) addCall(name string, args ...any) {
	db.calls[name] = append(db.calls[name], args)
}

// GetCalls returns calls for a given function name.
func (db *MockDatabaseAPI) GetCalls(name string) [][]any {
	return db.calls[name]
}

func (db *MockDatabaseAPI) GetUser(ctx context.Context, id uuid.UUID) (models.User, error) {
	db.addCall("GetUser", id)

	return db.GetUserResult.Get()
}

func (db *MockDatabaseAPI) GetUserWithCredentials(
	ctx context.Context,
	username string,
	password string,
) (models.User, error) {
	db.addCall("GetUserWithCredentials", username, password)

	return db.GetUserWithCredentialsResult.Get()
}

func (db *MockDatabaseAPI) RegisterUser(ctx context.Context, user models.User, password string) error {
	db.addCall("RegisterUser", user, password)

	return db.RegisterUserResult
}

func (db *MockDatabaseAPI) GetLanguage(ctx context.Context, id string) (models.Language, error) {
	db.addCall("GetLanguage", id)

	return db.GetLanguageResult.Get()
}

func (db *MockDatabaseAPI) FindCodeSamples(
	ctx context.Context,
	search models.CodeSampleSearch,
) (models.CodeSamplePage, error) {
	db.addCall("FindCodeSamples", search)

	return db.FindCodeSamplesResult.Get()
}

func (db *MockDatabaseAPI) GetCodeSample(ctx context.Context, id uuid.UUID) (models.CodeSample, error) {
	db.addCall("GetCodeSample", id)

	return db.GetCodeSampleResult.Get()
}

func (db *MockDatabaseAPI) CreateCodeSample(ctx context.Context, sample models.CodeSample) error {
	db.addCall("CreateCodeSample", sample)

	return db.CreateCodeSampleResult
}

func (db *MockDatabaseAPI) UpdateCodeSample(ctx context.Context, sample models.CodeSample) error {
	db.addCall("UpdateCodeSample", sample)

	return db.UpdateCodeSampleResult
}

func (db *MockDatabaseAPI) DeleteCodeSample(ctx context.Context, id uuid.UUID) error {
	db.addCall("DeleteCodeSample", id)

	return db.DeleteCodeSampleResult
}

func New() *MockDatabaseAPI {
	return &MockDatabaseAPI{
		calls: make(map[string][][]any),
	}
}

// Ensure the mock conforms with the interface.
func dummy() {
	var _ database.DatabaseAPI = New()
}
