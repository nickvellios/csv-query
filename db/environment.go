package db

import (
	"os"
)

func setEnv() {
	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_PASS", "")
	os.Setenv("DB_NAME", "postgres")
	os.Setenv("DB_DRIVER", "postgres")
	os.Setenv("DB_TABLE", "csvtest")
}

func SetupEnvironment() *UrlDB {
	setEnv()

	udb := &UrlDB{}
	err := udb.Open()
	if err != nil {
		panic(err)
	}
	return udb
}
