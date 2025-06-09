package store

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"kws/kws/consts/status"
	"kws/kws/models"
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

func (auth *AuthStore) VerifyUser(ctx context.Context, email string) error {
	sql := `
		UPDATE users SET verified=true WHERE email = $1
	`

	_, err := auth.db.Exec(ctx, sql, email)
	if err != nil {
		return err
	}

	return nil
}

func (auth *AuthStore) LoginUser(ctx context.Context, userName, password string) (*models.User, error) {
	var userModel models.User

	sql := `
		SELECT verified FROM users WHERE user_name = $1
	`

	err := auth.db.QueryRow(ctx, sql, userName).Scan(&userModel.Verified)

	if err != nil {
		log.Println("No rows found in the table")
		return nil, errors.New(status.USER_NAME_INVALID)
	}

	if !userModel.Verified {
		log.Println("The user is still not verified")
		return nil, errors.New(status.USER_NOT_VERIFIED)
	}

	sql = `
		SELECT id, first_name, last_name, email, password_hash, verified, user_name FROM users where user_name = $1
	`

	err = auth.db.QueryRow(ctx, sql, userName).Scan(
		&userModel.Id,
		&userModel.First_name,
		&userModel.Last_name,
		&userModel.Email,
		&userModel.Password,
		&userModel.Verified,
		&userModel.User_name,
	)
	if err != nil {
		log.Println("No rows found in the table: ", err.Error())
		return nil, errors.New(status.USER_NAME_INVALID)
	}

	err = bcrypt.CompareHashAndPassword([]byte(userModel.Password), []byte(password))
	if err != nil {
		log.Println("Wrong password")
		return nil, errors.New(status.WRONG_CREDENTIALS)
	}

	return &userModel, nil
}
