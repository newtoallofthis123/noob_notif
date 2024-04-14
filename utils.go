package main

import (
	"os"

	"github.com/joho/godotenv"
)

func GetDbUrl() string {
	godotenv.Load()

	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		panic("Env variable DB_URL not set")
	}

	return dbUrl
}
