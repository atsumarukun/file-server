package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	PORT int
)

func Load() error {
	var err error
	if err = godotenv.Load(".env"); err != nil {
		return err
	}

	if PORT, err = strconv.Atoi(os.Getenv("PORT")); err != nil {
		return err
	}

	return nil
}
