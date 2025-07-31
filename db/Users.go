package db

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type userDB map[string]string

var Users userDB

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	Users = make(map[string]string)
	adminUsername := strings.Trim(os.Getenv("AdminUsername"), " \n\t")
	adminPassword := strings.Trim(os.Getenv("AdminPassword"), " \n\t")
	Users[adminUsername] = adminPassword
}
