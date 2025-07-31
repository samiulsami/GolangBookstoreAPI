package auth

import (
	"GoBookstoreAPI/db"
	"GoBookstoreAPI/opentelemetry"
	"GoBookstoreAPI/prometheus_metrics"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"go.opentelemetry.io/otel"
)

// /middleware
func BasicAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		tracer := otel.Tracer(opentelemetry.ServiceName)

		_, span := tracer.Start(ctx, "BasicAuth Endpoint")
		defer span.End()

		span.AddEvent("Fetching credentials from request")
		prometheus_metrics.BasicAuthAttempts.Inc()
		authHeader := req.Header.Get("Authorization")
		if authHeader == "" {
			span.RecordError(fmt.Errorf("failed to get authorization header"))
			res.WriteHeader(http.StatusUnauthorized)
			if _, err := res.Write([]byte("Authorization header not found")); err != nil {
				span.RecordError(fmt.Errorf("failed to write response: %w", err))
				fmt.Println(err)
			}
			return
		}

		authHeaderSplit := strings.Split(authHeader, " ")

		span.AddEvent("Parsing credentials")
		if len(authHeaderSplit) != 2 {
			span.AddEvent("Invalid authorization header")
			res.WriteHeader(http.StatusUnauthorized)
			if _, err := res.Write([]byte("Invalid authorization credentials")); err != nil {
				span.RecordError(fmt.Errorf("failed to write response: %w", err))
				fmt.Println(err)
			}
			return
		}

		authType, encodedCredentials := authHeaderSplit[0], authHeaderSplit[1]

		if authType != "Basic" {
			span.AddEvent("Auth type invalid")
			res.WriteHeader(http.StatusUnauthorized)
			if _, err := res.Write([]byte("Invalid authorization type")); err != nil {
				span.RecordError(fmt.Errorf("failed to write response: %w", err))
				fmt.Println(err)
			}
			return
		}

		span.AddEvent("Decoding credentials")
		decodedCredentials, err := base64.StdEncoding.DecodeString(encodedCredentials)
		if err != nil {
			span.RecordError(fmt.Errorf("failed to decode credentials: %w", err))
			res.WriteHeader(http.StatusUnauthorized)
			if _, err := res.Write([]byte("Failed to decode credentials")); err != nil {
				span.RecordError(fmt.Errorf("failed to write response: %w", err))
				fmt.Println(err)
			}
			return
		}

		decodedString := string(decodedCredentials)
		idx := strings.Index(decodedString, ":")

		username, password := strings.Trim(decodedString[:idx], " \n\t"), strings.Trim(decodedString[idx+1:], " \n\t")

		span.AddEvent("Verifying credentials")
		if pass, ok := db.Users[username]; !ok || pass != password {
			span.AddEvent("Invalid credentials")
			res.WriteHeader(http.StatusUnauthorized)
			if _, err := res.Write([]byte("Invalid Credentials")); err != nil {
				span.RecordError(fmt.Errorf("failed to write response: %w", err))
				fmt.Println(err)
			}
			return
		}

		span.AddEvent("Basic auth successful")
		next.ServeHTTP(res, req)
		prometheus_metrics.BasicAuthSuccess.Inc()
	})
}
