package database

import (
	"context"

	"github.com/dense-analysis/codelibrary/internal/api/models"
)

func (db *databaseAPIImpl) GetLanguage(ctx context.Context, id string) (models.Language, error) {
	row := db.pool.QueryRow(ctx, `SELECT name FROM language WHERE id = $1`, id)

	language := models.Language{ID: id}
	err := row.Scan(&language.Name)

	return language, err
}
