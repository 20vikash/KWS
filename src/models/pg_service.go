package models

type PGServiceUser struct {
	Uid      int
	UserName string
	Password string
}

type PGServiceDatabase struct {
	Pid    int
	DbName string
}

func CreatePgServiceUser(uid int, userName, password string) *PGServiceUser {
	return &PGServiceUser{
		Uid:      uid,
		UserName: userName,
		Password: password,
	}
}

func CreatePgServiceDatabase(pid int, dbName string) *PGServiceDatabase {
	return &PGServiceDatabase{
		Pid:    pid,
		DbName: dbName,
	}
}
