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

type JobResponseDeploy struct {
	Done     bool
	Instance Instance
}

type JobResponseSK struct {
	Done    bool
	Success bool
}

type Instance struct {
	ContainerID string
	Success     bool
	Username    string
	Password    string
	IP          string
}

type InsData struct {
	ContainerName  string
	Username       string
	InstanceStatus string // "inactive", "active", "stopped"
	Active         string // "exists", "no"
	Instance       Instance
}

type Domain struct {
	Name   string
	Port   int
	Status string
}

type PublishInstancePageData struct {
	LoggedInUser string
	InstanceName string
	Domains      []Domain
	HasDomains   bool
}
