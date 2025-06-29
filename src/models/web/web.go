package web

type User struct {
	Username    string
	Password    string
	ID          int
	Permissions string
}

type PGUserPageData struct {
	HostName       string
	ServiceIP      string
	Port           string
	ActiveInstance string
	LoggedInUser   string
	UserLimit      int
	Users          []User
}
