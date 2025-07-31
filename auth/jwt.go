package auth

import (
	"GoBookstoreAPI/opentelemetry"
	"GoBookstoreAPI/prometheus_metrics"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/golang-jwt/jwt/v5/request"
	"github.com/joho/godotenv"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

var (
	ValidityDurationInSeconds int32 = 300
	secret                    []byte
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	secret = []byte(strings.Trim(os.Getenv("jwtsecret"), " \n\t"))
}

// /Sets cookie and returns the JWT token
func GetJWTToken(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	tracer := otel.Tracer(opentelemetry.ServiceName)

	_, span := tracer.Start(ctx, "GetJWTToken Endpoint")
	defer span.End()

	span.AddEvent("Decoding request body")
	var bodyMap map[string]string
	err := json.NewDecoder(req.Body).Decode(&bodyMap)

	username := ""
	if err != nil {
		span.RecordError(fmt.Errorf("failed to parse body: %w", err))
		username = strings.Trim(os.Getenv("adminusername"), " \n\t")
	} else if name, ok := bodyMap["username"]; ok {
		span.AddEvent("username found in request body")
		username = name
	}

	span.SetAttributes(attribute.String("username", username))
	token, err := createToken(username)
	if err != nil {
		span.RecordError(fmt.Errorf("failed to create token: %w", err))
		res.WriteHeader(http.StatusInternalServerError)
		if _, err := res.Write([]byte("Unable to create bearer token")); err != nil {
			span.RecordError(fmt.Errorf("failed to write response: %w", err))
			fmt.Println(err)
		}
		return
	}

	cookie := http.Cookie{
		Name:    "jwt",
		Value:   token,
		Expires: time.Now().Add(time.Second * time.Duration(ValidityDurationInSeconds)),
		Path:    "/",
	}

	span.AddEvent("Setting cookie")
	http.SetCookie(res, &cookie)
	res.WriteHeader(http.StatusOK)
	if _, err := res.Write([]byte("Bearer Token: " + token + "\nCookie set")); err != nil {
		span.RecordError(fmt.Errorf("failed to write response: %w", err))
		fmt.Println(err)
	}
}

func JWTAuthenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		tracer := otel.Tracer(opentelemetry.ServiceName)

		_, span := tracer.Start(ctx, "JWTAuthenticator Endpoint")
		defer span.End()

		prometheus_metrics.JWTAuthAttempts.Inc()
		extractor := request.BearerExtractor{}
		span.AddEvent("Extracting bearer token")
		bearerToken, err := extractor.ExtractToken(req)

		switch err {
		case nil:
			span.AddEvent("Bearer token extracted")
			span.AddEvent("Verifying token")
			valid, err := verifyToken(bearerToken)
			if err != nil || !valid {
				span.RecordError(fmt.Errorf("failed to verify token: %w", err))
				res.WriteHeader(http.StatusUnauthorized)
				if _, err := res.Write([]byte("Invalid bearer token")); err != nil {
					span.RecordError(fmt.Errorf("failed to write response: %w", err))
					fmt.Println(err)
				}
				return
			}

			res.WriteHeader(http.StatusOK)
			next.ServeHTTP(res, req)

		default:
			span.AddEvent("Failed to extract bearer token")
			span.AddEvent("Attempting to fetch token from cookie")
			cookie, err := req.Cookie("jwt")
			if err != nil {
				span.RecordError(fmt.Errorf("failed to get cookie: %w", err))
				res.WriteHeader(http.StatusBadRequest)
				if _, err := res.Write([]byte("Cookie is invalid or nonexistent")); err != nil {
					span.RecordError(fmt.Errorf("failed to write response: %w", err))
					fmt.Println(err)
				}
				return
			}

			span.AddEvent("Cookie found")
			span.AddEvent("Verifying token")
			token := cookie.Value
			valid, err := verifyToken(token)
			if err != nil {
				span.RecordError(fmt.Errorf("failed to verify token: %w", err))
				res.WriteHeader(http.StatusBadRequest)
				fmt.Println(err)
				return
			}

			if !valid {
				span.AddEvent("Token invalid")
				res.WriteHeader(http.StatusUnauthorized)
				if _, err := res.Write([]byte("Authorization failed")); err != nil {
					span.RecordError(fmt.Errorf("failed to write response: %w", err))
					fmt.Println(err)
				}
				return
			}

			span.AddEvent("Token validation successful")
			next.ServeHTTP(res, req)
			prometheus_metrics.JWTAuthSuccess.Inc()
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
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
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
