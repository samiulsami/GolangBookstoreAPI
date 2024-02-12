package auth

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/golang-jwt/jwt/v5/request"
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
	var bodyMap map[string]string
	err := json.NewDecoder(req.Body).Decode(&bodyMap)

	var username string = ""
	if err != nil {
		username = "DefaultUsername"
	} else if name, ok := bodyMap["username"]; ok {
		username = name
	} else {
		username = os.Getenv("username")
	}

	token, err := createToken(username)

	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte("Unable to create bearer token"))
	}

	cookie := http.Cookie{
		Name:    "jwt",
		Value:   token,
		Expires: time.Now().Add(time.Second * time.Duration(ValidityDurationInSeconds)),
		Path:    "/",
	}

	http.SetCookie(res, &cookie)
	res.WriteHeader(http.StatusOK)
	res.Write([]byte("Bearer Token: " + token + "\nCookie set"))
}

func JWTAuthenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {

		extractor := request.BearerExtractor{}
		bearerToken, err := extractor.ExtractToken(req)

		switch err {
		case nil:
			valid, err := verifyToken(bearerToken)
			if err != nil || !valid {
				res.WriteHeader(http.StatusUnauthorized)
				res.Write([]byte("Invalid bearer token"))
				return
			}

			res.WriteHeader(http.StatusOK)
			next.ServeHTTP(res, req)

		default:
			cookie, err := req.Cookie("jwt")

			if err != nil {
				res.WriteHeader(http.StatusBadRequest)
				res.Write([]byte("Cookie is invalid or nonexistent"))
				return
			}

			token := cookie.Value
			valid, err := verifyToken(token)

			if err != nil {
				res.WriteHeader(http.StatusBadRequest)
				fmt.Println(err)
				return
			}

			if !valid {
				res.WriteHeader(http.StatusUnauthorized)
				res.Write([]byte("Authorization failed"))
				return
			}
			next.ServeHTTP(res, req)
		}
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
	fmt.Printf("JWT Token validity period: %d hours %d minutes and %d seconds\n", remainingSeconds/3600, (remainingSeconds%3600)/60, remainingSeconds%60)

	return true, nil
}
