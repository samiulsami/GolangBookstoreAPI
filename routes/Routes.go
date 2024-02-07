package routes

import (
	"GoBookstoreAPI/auth"
	"GoBookstoreAPI/handlers"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func StartAPI(host string, port string) {
	r := chi.NewRouter()
	r.Use(CommonHeaders)

	r.Group(func(r chi.Router) {
		r.Use(auth.BasicAuth)
		r.Get("/api/v1/get-token", auth.GetJWTToken)
	})

	r.Group(func(r chi.Router) {
		r.Use(auth.JWTAuthenticator)
		r.Post("/api/v1/books", handlers.AddBook)
		r.Get("/api/v1/books", handlers.GetAllBooks)
		r.Get("/api/v1/books/{id}", handlers.GetBook)
		r.Delete("/api/v1/books/{id}", handlers.DeleteBook)
		r.Put("/api/v1/books/{id}", handlers.UpdateBook)
	})

	http.ListenAndServe(host+":"+port, r)
}

func CommonHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(res, req)
	})
}
