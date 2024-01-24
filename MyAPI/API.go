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
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(Auth.BasicAuth)
		r.Get("/api/v1/get-token", Auth.GetJWTToken)
	})

	http.ListenAndServe("localhost:"+port, r)
}
