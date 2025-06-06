package models

import "regexp"

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
