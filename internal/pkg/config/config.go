package config

import (
	"fmt"
	"os"
	"strconv"
)

const (
	STORAGE_PATH = "storage"
)

var (
	API_PORT       int
	MYSQL_DSN      string
	JWT_SECRET_KEY string
)

func Load() error {
	var err error
	if API_PORT, err = strconv.Atoi(os.Getenv("API_PORT")); err != nil {
		return err
	}

	var databasePort int
	if databasePort, err = strconv.Atoi(os.Getenv("MYSQL_PORT")); err != nil {
		return err
	}
	MYSQL_DSN = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", os.Getenv("MYSQL_USER"), os.Getenv("MYSQL_PASSWORD"), os.Getenv("MYSQL_HOST"), databasePort, os.Getenv("MYSQL_DATABASE"))

	JWT_SECRET_KEY = os.Getenv("JWT_SECRET_KEY")

	return nil
}
