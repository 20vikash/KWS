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

type Database struct {
	Name  string
	Owner string
}

type PgDatabase struct {
	Password       string
	HostName       string
	Username       string
	TotalDatabases int
	Owner          string
	AvailableSlots int
	Limit          int
	UsagePercent   int
	Databases      []Database
}
