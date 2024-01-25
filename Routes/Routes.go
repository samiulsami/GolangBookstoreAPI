package Routes

import (
	"GoBookstoreAPI/Auth"
	"GoBookstoreAPI/Handlers"
	"github.com/go-chi/chi/v5"
	"net/http"
)

var port string = "3000"

func StartAPI() {
	r := chi.NewRouter()
	r.Use(CommonHeaders)

	r.Group(func(r chi.Router) {
		r.Use(Auth.BasicAuth)
		r.Get("/api/v1/get-token", Auth.GetJWTToken)
	})

	r.Group(func(r chi.Router) {
		r.Use(Auth.JWTAuthenticator)
		r.Post("/api/v1/books", Handlers.AddBook)
		r.Get("/api/v1/books", Handlers.GetAllBooks)
		r.Get("/api/v1/books/{id}", Handlers.GetBook)
		r.Delete("/api/v1/books/{id}", Handlers.DeleteBook)
		r.Put("/api/v1/books/{id}", Handlers.UpdateBook)
	})

	http.ListenAndServe("localhost:"+port, r)
}

func CommonHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(res, req)
	})
}
