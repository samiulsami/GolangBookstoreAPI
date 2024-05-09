package db

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strings"
)

type userDB map[string]string

var Users userDB

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	Users = make(map[string]string)
	adminUsername := strings.Trim(os.Getenv("AdminUsername"), " \n\n\t")
	adminPassword := strings.Trim(os.Getenv("AdminPassword"), " \n\n\t")
	Users[adminUsername] = adminPassword
}
