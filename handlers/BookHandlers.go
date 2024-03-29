package handlers

import (
	"GoBookstoreAPI/db"
	"GoBookstoreAPI/prometheusMetrics"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func AddBook(res http.ResponseWriter, req *http.Request) {
	var newBook db.Book
	err := json.NewDecoder(req.Body).Decode(&newBook)

	if err != nil || !newBook.IsValid() {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("Failed to parse body. Invalid book format"))
		return
	}

	uuid := db.BookDB.AddBook(&newBook)
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte("Book added. UUID: " + uuid))
	prometheusMetrics.BookAddCounter.Inc()
}

func GetBook(res http.ResponseWriter, req *http.Request) {
	uuid := chi.URLParam(req, "id")
	body, err := db.BookDB.GetBook(uuid)

	if err != nil {
		fmt.Println(err)
		res.WriteHeader(http.StatusNotFound)
		return
	}

	res.WriteHeader(http.StatusOK)
	res.Write(body)
	prometheusMetrics.BookGetCounter.Inc()
}

func GetAllBooks(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusOK)
	res.Write(db.BookDB.GetBookList())
	prometheusMetrics.BookGetAllCounter.Inc()
}

func DeleteBook(res http.ResponseWriter, req *http.Request) {
	uuid := chi.URLParam(req, "id")
	done, err := db.BookDB.DeleteBook(uuid)

	if err != nil || !done {
		fmt.Println(err)
		res.WriteHeader(http.StatusNotFound)
		return
	}

	res.WriteHeader(http.StatusOK)
	res.Write([]byte("Successfully deleted the book with UUID: " + uuid))
	prometheusMetrics.BookDeleteCounter.Inc()
}

func UpdateBook(res http.ResponseWriter, req *http.Request) {
	uuid := chi.URLParam(req, "id")
	var newBook db.Book
	err := json.NewDecoder(req.Body).Decode(&newBook)

	if err != nil || !newBook.IsValid() {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("Failed to parse body. Invalid book format"))
		return
	}

	newBook.UUID = uuid
	done, err := db.BookDB.UpdateBook(&newBook)

	if err != nil || !done {
		res.WriteHeader(http.StatusNotFound)
		res.Write([]byte("Book not found"))
		return
	}

	res.WriteHeader(http.StatusCreated)
	res.Write([]byte("Successfully updated the book with UUID: " + uuid))
	prometheusMetrics.BookUpdateCounter.Inc()
}
