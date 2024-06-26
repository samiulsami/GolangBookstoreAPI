package auth

import (
	"GoBookstoreAPI/db"
	"GoBookstoreAPI/prometheusMetrics"
	"encoding/base64"
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
			res.Write([]byte("Authorization header not found"))
			return
		}

		authHeaderSplit := strings.Split(authHeader, " ")

		if len(authHeaderSplit) != 2 {
			res.WriteHeader(http.StatusUnauthorized)
			res.Write([]byte("Invalid authorization credentials"))
			return
		}

		authType, encodedCredentials := authHeaderSplit[0], authHeaderSplit[1]

		if authType != "Basic" {
			res.WriteHeader(http.StatusUnauthorized)
			res.Write([]byte("Authorization type must be \"Basic\""))
			return
		}

		decodedCredentials, err := base64.StdEncoding.DecodeString(encodedCredentials)

		if err != nil {
			res.WriteHeader(http.StatusUnauthorized)
			res.Write([]byte("Failed to decode credentials"))
			return
		}

		decodedString := string(decodedCredentials)
		idx := strings.Index(decodedString, ":")

		username, password := strings.Trim(decodedString[:idx], " \n\n\t"), strings.Trim(decodedString[idx+1:], " \n\n\t")

		if pass, ok := db.Users[username]; !ok || pass != password {
			res.WriteHeader(http.StatusUnauthorized)
			res.Write([]byte("Invalid Credentials"))
			return
		}

		next.ServeHTTP(res, req)
		prometheusMetrics.BasicAuthSuccess.Inc()
	})
}
