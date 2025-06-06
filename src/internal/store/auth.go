package store

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

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

func (auth *AuthStore) GenerateToken(ctx context.Context, email string) string {
	uuid, err := exec.Command("uuidgen").Output()
	if err != nil {
		log.Fatal(err)
	}

	t := email + time.Now().String() + fmt.Sprintf("%x", uuid)

	h := sha256.New()
	h.Write([]byte(t))
	bs := h.Sum(nil)

	token := fmt.Sprintf("%x", bs)

	return token
}

func (auth *AuthStore) CreateUser(ctx context.Context, first_name, last_name, email, password, user_name string) error {
	passwordHash := hashedPassword(password)

	sql := `
		INSERT INTO users(first_name, last_name, email, password_hash, user_name) VALUES($1, $2, $3, $4, $5)
	`

	_, err := auth.db.Exec(ctx, sql,
		first_name,
		last_name,
		email,
		passwordHash,
		user_name,
	)
	if err != nil {
		if strings.Contains(err.Error(), "23505") { // 23505 is the error code for duplicate violation
			log.Println("Error: Duplicate entry (unique constraint violation)")
			return errors.New("username or email already exists")
		} else {
			log.Println("Cannot create a new User. Failed at inserting the user details into the table.")
		}

		return err
	}

	return nil
}
