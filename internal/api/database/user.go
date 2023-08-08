package database

import (
	"context"

	"github.com/dense-analysis/codelibrary/internal/api/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (db *databaseAPIImpl) GetUser(ctx context.Context, id uuid.UUID) (models.User, error) {
	row := db.pool.QueryRow(ctx, `SELECT username FROM "user" WHERE id = $1`, id)

	user := models.User{ID: id}
	err := row.Scan(&user.Username)

	return user, err
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (db *databaseAPIImpl) GetUserWithCredentials(
	ctx context.Context,
	username string,
	password string,
) (models.User, error) {
	row := db.pool.QueryRow(ctx, `SELECT id, password_hash FROM "user" WHERE username = $1`, username)

	user := models.User{Username: username}
	var hash string
	err := row.Scan(&user.ID, &hash)

	if err == nil && !CheckPasswordHash(password, hash) {
		err = NotFoundErr
	}

	return user, err
}

func (db *databaseAPIImpl) RegisterUser(ctx context.Context, user models.User, password string) error {
	hash, err := HashPassword(password)

	if err != nil {
		return err
	}

	_, err = db.pool.Exec(
		ctx,
		`INSERT INTO "user" (id, username, password_hash) VALUES ($1, $2, $3)`,
		user.ID, user.Username, hash,
	)

	return err
}
