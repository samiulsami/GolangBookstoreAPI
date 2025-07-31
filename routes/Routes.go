package routes

import (
	"GoBookstoreAPI/auth"
	"GoBookstoreAPI/handlers"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func StartAPI(host string, port string) {
	r := chi.NewRouter()
	r.Use(CommonHeaders)
	r.Handle("/metrics", promhttp.Handler())

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

	if err := http.ListenAndServe(host+":"+port, r); err != nil {
		log.Fatal(err)
	}
}

func CommonHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(res, req)
	})
}
