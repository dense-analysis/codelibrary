package database

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/dense-analysis/codelibrary/internal/api/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var NotFoundErr = pgx.ErrNoRows
var DuplicateErr = errors.New("duplicate object")

type DatabaseAPI interface {
	GetUser(ctx context.Context, id uuid.UUID) (models.User, error)
	GetUserWithCredentials(ctx context.Context, username string, password string) (models.User, error)
	RegisterUser(ctx context.Context, user models.User, password string) error
	GetLanguage(ctx context.Context, id string) (models.Language, error)
	FindCodeSamples(ctx context.Context, search models.CodeSampleSearch) (models.CodeSamplePage, error)
	GetCodeSample(ctx context.Context, id uuid.UUID) (models.CodeSample, error)
	CreateCodeSample(ctx context.Context, sample models.CodeSample) error
	UpdateCodeSample(ctx context.Context, sample models.CodeSample) error
	DeleteCodeSample(ctx context.Context, id uuid.UUID) error
}

type ConnectionPool interface {
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
	Begin(context.Context) (pgx.Tx, error)
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error)
	Close()
}

type databaseAPIImpl struct {
	pool ConnectionPool
}

func NewWithPool(pool ConnectionPool) (DatabaseAPI, error) {
	return &databaseAPIImpl{pool: pool}, nil
}

func New(ctx context.Context) (DatabaseAPI, error) {
	config, err := pgxpool.ParseConfig(fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
	))

	if err != nil {
		return nil, nil
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)

	if err != nil {
		return nil, err
	}

	return NewWithPool(pool)
}
