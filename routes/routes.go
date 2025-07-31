package routes

import (
	"GoBookstoreAPI/auth"
	"GoBookstoreAPI/handlers"
	"GoBookstoreAPI/opentelemetry"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func StartAPI(host string, port string) {
	shutdown, err := opentelemetry.InitTracer()
	if err != nil {
		log.Fatalf("error initializing tracer provider: %v", err)
	}
	defer shutdown()

	r := chi.NewRouter()
	r.Use(CommonHeaders)
	r.Handle("/metrics", promhttp.Handler())

	r.With(
		func(next http.Handler) http.Handler {
			return otelhttp.NewHandler(next, "bookstore-token-http-request")
		}).Group(func(r chi.Router) {
		r.Use(auth.BasicAuth)
		r.Get("/api/v1/get-token", auth.GetJWTToken)
	})

	r.With(func(next http.Handler) http.Handler {
		return otelhttp.NewHandler(next, "bookstore-http-request")
	}).Group(func(r chi.Router) {
		r.Use(auth.JWTAuthenticator)
		r.Post("/api/v1/books", handlers.AddBook)
		r.Get("/api/v1/books", handlers.GetAllBooks)
		r.Get("/api/v1/books/{id}", handlers.GetBook)
		r.Delete("/api/v1/books/{id}", handlers.DeleteBook)
		r.Put("/api/v1/books/{id}", handlers.UpdateBook)
	})

	if err := http.ListenAndServe(host+":"+port, r); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

func CommonHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(res, req)
	})
}
