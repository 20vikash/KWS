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
}

func (u *User) ValidateEmail() bool {
	regex := `(?i)^[a-zA-Z0-9._%+\-]+@gmail\.com$`

	re := regexp.MustCompile(regex)
	return re.MatchString(u.Email)
}

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
