package Auth

import (
	"GoBookstoreAPI/DB"
	"encoding/base64"
	"net/http"
	"strings"
)

// /middleware
func BasicAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		authHeader := req.Header.Get("Authorization")
		if authHeader == "" {
			res.Write([]byte("Authorization header not found"))
			res.WriteHeader(http.StatusUnauthorized)
			return
		}

		authHeaderSplit := strings.Split(authHeader, " ")

		if len(authHeaderSplit) != 2 {
			res.Write([]byte("Invalid authorization credentials"))
			res.WriteHeader(http.StatusUnauthorized)
			return
		}

		authType, encodedCredentials := authHeaderSplit[0], authHeaderSplit[1]

		if authType != "Basic" {
			res.Write([]byte("Authorization type must be \"Basic\""))
			res.WriteHeader(http.StatusUnauthorized)
			return
		}

		decodedCredentials, err := base64.StdEncoding.DecodeString(encodedCredentials)

		if err != nil {
			res.Write([]byte("Failed to decode credentials"))
			res.WriteHeader(http.StatusUnauthorized)
			return
		}

		decodedString := string(decodedCredentials)
		idx := strings.Index(decodedString, ":")

		username, password := decodedString[:idx], decodedString[idx+1:]

		if pass, ok := DB.Users[username]; !ok || pass != password {
			res.Write([]byte("Invalid Credentials"))
			res.WriteHeader(http.StatusUnauthorized)
			return
		}

		req.Header.Add("username", username)
		next.ServeHTTP(res, req)
	})
}
