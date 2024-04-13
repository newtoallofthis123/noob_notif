package main

import (
	"math/rand"
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

func generate(length int, pool string) string {
	var otp string

	for i := 0; i < length; i++ {
		otp += string(pool[rand.Intn(len(pool))])
	}

	return otp
}

func RanHash(len int) string {
	pool := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	return generate(len, pool)
}
