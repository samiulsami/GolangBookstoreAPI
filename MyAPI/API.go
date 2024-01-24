package MyAPI

import (
	"GoBookstoreAPI/Auth"
	"GoBookstoreAPI/DB"
	"github.com/go-chi/chi/v5"
	"net/http"
)

var port string = "3000"

func StartAPI() {
	DB.Init()
	Auth.Init()

	r := chi.NewRouter()
	r.Use(CommonHeaders)

	r.Group(func(r chi.Router) {
		r.Use(Auth.BasicAuth)
		r.Get("/api/v1/get-token", Auth.GetJWTToken)
	})

	http.ListenAndServe("localhost:"+port, r)
}

func CommonHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(res, req)
	})
}
