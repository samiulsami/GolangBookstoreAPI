package Auth

import (
	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"net/http"
	"time"
)

var ValidityDurationInSeconds int32 = 300
var secret string = "orangeCat"
var tokenAuth *jwtauth.JWTAuth

func Init() {
	tokenAuth = jwtauth.New("HS256", []byte(secret), nil, jwt.WithAcceptableSkew(30*time.Second))
}

// /Sets cookie and returns the JWT token
func GetJWTToken(res http.ResponseWriter, req *http.Request) {
	username := req.Header.Get("username")
	if username == "" {
		res.Write([]byte("This should be impossible"))
		res.WriteHeader(http.StatusUnauthorized)
		return
	}

	token := createToken(username)

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

func createToken(username string) string {
	exp := time.Now().Add(time.Second * time.Duration(ValidityDurationInSeconds)).Unix()
	iat := time.Now().Unix()

	_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{
		"sub": username,
		"exp": exp,
		"iat": iat,
	})

	return tokenString
}
