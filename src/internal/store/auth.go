package store

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type AuthStore struct {
	db *pgxpool.Pool
}

func hashedPassword(password string) []byte {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Failed to hash the password")
	}

	return hash
}

func (auth *AuthStore) CreateUser(ctx context.Context, first_name, last_name, email, password string) error {
	passwordHash := hashedPassword(password)

	sql := `
		INSERT INTO users(first_name, last_name, email, password_hash) VALUES($1, $2, %3. $4)
	`

	_, err := auth.db.Exec(ctx, sql,
		first_name,
		last_name,
		email,
		passwordHash,
	)
	if err != nil {
		log.Println("Cannot create a new User. Failed at inserting the user details into the table.")
		return err
	}

	return nil
}
