package web

type User struct {
	Role        string
	Username    string
	Password    string
	ID          int
	Permissions string
}

type PGUserPageData struct {
	HostName     string
	ServiceIP    string
	Port         string
	LoggedInUser string
	UserLimit    int
	Users        []User
}
