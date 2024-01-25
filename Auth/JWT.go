package Auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"time"
)

var ValidityDurationInSeconds int32 = 300
var secret []byte

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	secret = []byte(os.Getenv("JWTSECRET"))
}

// /Sets cookie and returns the JWT token
func GetJWTToken(res http.ResponseWriter, req *http.Request) {
	username := req.Header.Get("username")
	if username == "" {
		res.Write([]byte("This should be impossible"))
		res.WriteHeader(http.StatusUnauthorized)
		return
	}

	token, err := createToken(username)

	if err != nil {
		res.Write([]byte("Unable to create bearer token"))
		res.WriteHeader(http.StatusInternalServerError)
	}

	cookie := http.Cookie{
		Name:    "jwt",
		Value:   token,
		Expires: time.Now().Add(time.Second * time.Duration(ValidityDurationInSeconds)),
		Path:    "/",
	}

	http.SetCookie(res, &cookie)
	res.Write([]byte("Bearer Token: " + token + "\nCookie set"))
	res.WriteHeader(http.StatusOK)
}

func JWTAuthenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		cookie, err := req.Cookie("jwt")

		if err != nil {
			res.Write([]byte("Cookie is invalid or nonexistent"))
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		token := cookie.Value
		valid, err := verifyToken(token)

		if err != nil {
			fmt.Println(err)
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		if !valid {
			res.Write([]byte("Authorization failed"))
			res.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(res, req)
	})
}

func createToken(username string) (string, error) {
	exp := time.Now().Add(time.Second * time.Duration(ValidityDurationInSeconds)).Unix()
	iat := time.Now().Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": username,
		"exp": exp,
		"iat": iat,
	})

	tokenString, err := token.SignedString(secret)
	return tokenString, err
}

func verifyToken(tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})
	if err != nil {
		return false, err
	}

	expInterface, ok := token.Claims.(jwt.MapClaims)["exp"]
	if !ok {
		return false, nil
	}

	remainingSeconds := (int64(expInterface.(float64)) - time.Now().Unix())
	if remainingSeconds <= 0 {
		return false, nil
	}
	fmt.Printf("Remaining time: %d hours %d minutes and %d seconds\n", remainingSeconds/3600, (remainingSeconds%3600)/60, remainingSeconds%60)

	return true, nil
}
