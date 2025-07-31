package auth

import (
	"GoBookstoreAPI/db"
	"GoBookstoreAPI/prometheusMetrics"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

// /middleware
func BasicAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		prometheusMetrics.BasicAuthAttempts.Inc()
		authHeader := req.Header.Get("Authorization")
		if authHeader == "" {
			res.WriteHeader(http.StatusUnauthorized)
			if _, err := res.Write([]byte("Authorization header not found")); err != nil {
				fmt.Println(err)
			}
			return
		}

		authHeaderSplit := strings.Split(authHeader, " ")

		if len(authHeaderSplit) != 2 {
			res.WriteHeader(http.StatusUnauthorized)
			if _, err := res.Write([]byte("Invalid authorization credentials")); err != nil {
				fmt.Println(err)
			}
			return
		}

		authType, encodedCredentials := authHeaderSplit[0], authHeaderSplit[1]

		if authType != "Basic" {
			res.WriteHeader(http.StatusUnauthorized)
			if _, err := res.Write([]byte("Invalid authorization type")); err != nil {
				fmt.Println(err)
			}
			return
		}

		decodedCredentials, err := base64.StdEncoding.DecodeString(encodedCredentials)
		if err != nil {
			res.WriteHeader(http.StatusUnauthorized)
			if _, err := res.Write([]byte("Failed to decode credentials")); err != nil {
				fmt.Println(err)
			}
			return
		}

		decodedString := string(decodedCredentials)
		idx := strings.Index(decodedString, ":")

		username, password := strings.Trim(decodedString[:idx], " \n\t"), strings.Trim(decodedString[idx+1:], " \n\t")

		if pass, ok := db.Users[username]; !ok || pass != password {
			res.WriteHeader(http.StatusUnauthorized)
			if _, err := res.Write([]byte("Invalid Credentials")); err != nil {
				fmt.Println(err)
			}
			return
		}

		next.ServeHTTP(res, req)
		prometheusMetrics.BasicAuthSuccess.Inc()
	})
}
