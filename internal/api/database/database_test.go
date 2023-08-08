package database_test

import (
	"testing"
	"time"

	"github.com/dense-analysis/codelibrary/internal/api/database"
	"github.com/pashagolub/pgxmock/v2"
)

type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v interface{}) bool {
	_, ok := v.(time.Time)
	return ok
}

func startDatabaseTest(t *testing.T) (pgxmock.PgxPoolIface, database.DatabaseAPI) {
	t.Parallel()
	t.Helper()

	mock, err := pgxmock.NewPool()

	if err != nil {
		t.Fatal(err)
	}

	db, err := database.NewWithPool(mock)

	if err != nil {
		t.Fatal(err)
	}

	return mock, db
}
