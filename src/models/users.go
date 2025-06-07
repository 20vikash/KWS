package models

import (
	"regexp"
	"strings"
)

type User struct {
	Email      string
	First_name string
	Last_name  string
	User_name  string
	Password   string
	Verified   bool
}

func (u *User) ValidateEmail() bool {
	regex := `(?i)^[a-zA-Z0-9._%+\-]+@gmail\.com$`

	re := regexp.MustCompile(regex)
	return re.MatchString(u.Email)
}

// Should be atleast 8 characters with atleast 1 low case, upper case, number, and a special character
func (u *User) ValidatePassword() bool {
	if len(u.Password) < 8 {
		return false
	}

	hasUpper := false
	hasLower := false
	hasNumber := false
	hasSpecial := false

	for _, c := range u.Password {
		switch {
		case 'a' <= c && 'z' >= c:
			hasLower = true
		case 'A' <= c && 'Z' >= c:
			hasUpper = true
		case '0' <= c && '9' >= c:
			hasNumber = true
		case strings.ContainsRune("!@#$%^&*()-_=+[]{}|;:'\",.<>/?", c):
			hasSpecial = true
		}
	}

	return (hasUpper && hasLower && hasNumber && hasSpecial)
}

// Should be from 5 to 30 characters. Should not start with a special character. Limited special characters.
func (u *User) ValidateUserName() bool {
	re := regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9._]{4,29}$`)
	return re.MatchString(u.User_name)
}

// Min 2 characters, max 30 characters. With limited special characters
func (u *User) ValidateFLNames() bool {
	re := regexp.MustCompile(`^[A-Za-z'.-]{2,30}$`)
	return re.MatchString(u.First_name) && re.MatchString(u.Last_name)
}
