package handlers

import (
	"GoBookstoreAPI/db"
	"GoBookstoreAPI/opentelemetry"
	"GoBookstoreAPI/prometheus_metrics"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/otel"
)

func AddBook(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	tracer := otel.Tracer(opentelemetry.ServiceName)

	_, span := tracer.Start(ctx, "AddBook Endpoint")
	defer span.End()

	span.AddEvent("Decoding request body")
	var newBook db.Book
	err := json.NewDecoder(req.Body).Decode(&newBook)
	if err != nil || !newBook.IsValid() {
		span.RecordError(fmt.Errorf("failed to parse body: %w", err))
		res.WriteHeader(http.StatusBadRequest)
		if _, err := res.Write([]byte("Failed to parse body. Invalid book format")); err != nil {
			span.RecordError(fmt.Errorf("failed to write response: %w", err))
			fmt.Println(err)
		}
		return
	}

	span.AddEvent("Adding book to database")
	uuid := db.BookDB.AddBook(&newBook)
	res.WriteHeader(http.StatusCreated)
	if _, err := res.Write([]byte("Book added. UUID: " + uuid)); err != nil {
		span.RecordError(fmt.Errorf("failed to write response: %w", err))
		fmt.Println(err)
	}
	prometheus_metrics.BookAddCounter.Inc()
}

func GetBook(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	tracer := otel.Tracer(opentelemetry.ServiceName)

	_, span := tracer.Start(ctx, "GetBook Endpoint")
	defer span.End()

	uuid := chi.URLParam(req, "id")
	span.AddEvent("Getting book from database")
	body, err := db.BookDB.GetBook(uuid)
	if err != nil {
		span.RecordError(fmt.Errorf("failed to get book from database: %w", err))
		fmt.Println(err)
		res.WriteHeader(http.StatusNotFound)
		return
	}

	span.AddEvent("Sending book to client")
	res.WriteHeader(http.StatusOK)
	if _, err := res.Write(body); err != nil {
		span.RecordError(fmt.Errorf("failed to write response: %w", err))
		fmt.Println(err)
	}
	prometheus_metrics.BookGetCounter.Inc()
}

func GetAllBooks(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	tracer := otel.Tracer(opentelemetry.ServiceName)

	_, span := tracer.Start(ctx, "GetAllBooks Endpoint")
	defer span.End()

	span.AddEvent("Getting all books from database")
	res.WriteHeader(http.StatusOK)
	if _, err := res.Write(db.BookDB.GetBookList()); err != nil {
		span.RecordError(fmt.Errorf("failed to write response: %w", err))
		fmt.Println(err)
	}
	prometheus_metrics.BookGetAllCounter.Inc()
}

func DeleteBook(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	tracer := otel.Tracer(opentelemetry.ServiceName)

	_, span := tracer.Start(ctx, "DeleteAllBooks Endpoint")
	defer span.End()

	span.AddEvent("Getting book from database")
	uuid := chi.URLParam(req, "id")
	done, err := db.BookDB.DeleteBook(uuid)

	if err != nil || !done {
		span.RecordError(fmt.Errorf("failed to delete book from database: %w", err))
		fmt.Println(err)
		res.WriteHeader(http.StatusNotFound)
		return
	}

	span.AddEvent("Sending response to client")
	res.WriteHeader(http.StatusOK)
	if _, err := res.Write([]byte("Successfully deleted the book with UUID: " + uuid)); err != nil {
		span.RecordError(fmt.Errorf("failed to write response: %w", err))
		fmt.Println(err)
	}
	prometheus_metrics.BookDeleteCounter.Inc()
}

func UpdateBook(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	tracer := otel.Tracer(opentelemetry.ServiceName)

	_, span := tracer.Start(ctx, "UpdateBook Endpoint")
	defer span.End()

	span.AddEvent("Decoding request body")
	uuid := chi.URLParam(req, "id")
	var newBook db.Book
	err := json.NewDecoder(req.Body).Decode(&newBook)
	if err != nil || !newBook.IsValid() {
		span.RecordError(fmt.Errorf("failed to parse body: %w", err))
		res.WriteHeader(http.StatusBadRequest)
		if _, err := res.Write([]byte("Failed to parse body. Invalid book format")); err != nil {
			span.RecordError(fmt.Errorf("failed to write response: %w", err))
			fmt.Println(err)
		}
		return
	}

	span.AddEvent("Updating book in database")
	newBook.UUID = uuid
	done, err := db.BookDB.UpdateBook(&newBook)

	if err != nil || !done {
		span.AddEvent("error updating book in database")
		res.WriteHeader(http.StatusNotFound)
		if _, err := res.Write([]byte("Book not found")); err != nil {
			span.RecordError(fmt.Errorf("failed to write response: %w", err))
			fmt.Println(err)
		}
		return
	}

	res.WriteHeader(http.StatusCreated)
	if _, err := res.Write([]byte("Successfully updated the book with UUID: " + uuid)); err != nil {
		span.RecordError(fmt.Errorf("failed to write response: %w", err))
		fmt.Println(err)
	}
	prometheus_metrics.BookUpdateCounter.Inc()
}
