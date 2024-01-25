package DB

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type userDB map[string]string

var Users userDB

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	Users = make(map[string]string)

	Users[os.Getenv("AdminUsername")] = os.Getenv("AdminPassword")
}
