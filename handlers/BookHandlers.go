package handlers

import (
	"GoBookstoreAPI/db"
	"GoBookstoreAPI/prometheusMetrics"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func AddBook(res http.ResponseWriter, req *http.Request) {
	var newBook db.Book
	err := json.NewDecoder(req.Body).Decode(&newBook)

	if err != nil || !newBook.IsValid() {
		res.WriteHeader(http.StatusBadRequest)
		if _, err := res.Write([]byte("Failed to parse body. Invalid book format")); err != nil {
			fmt.Println(err)
		}
		return
	}

	uuid := db.BookDB.AddBook(&newBook)
	res.WriteHeader(http.StatusCreated)
	if _, err := res.Write([]byte("Book added. UUID: " + uuid)); err != nil {
		fmt.Println(err)
	}
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
	if _, err := res.Write(body); err != nil {
		fmt.Println(err)
	}
	prometheusMetrics.BookGetCounter.Inc()
}

func GetAllBooks(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusOK)
	if _, err := res.Write(db.BookDB.GetBookList()); err != nil {
		fmt.Println(err)
	}
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
	if _, err := res.Write([]byte("Successfully deleted the book with UUID: " + uuid)); err != nil {
		fmt.Println(err)
	}
	prometheusMetrics.BookDeleteCounter.Inc()
}

func UpdateBook(res http.ResponseWriter, req *http.Request) {
	uuid := chi.URLParam(req, "id")
	var newBook db.Book
	err := json.NewDecoder(req.Body).Decode(&newBook)

	if err != nil || !newBook.IsValid() {
		res.WriteHeader(http.StatusBadRequest)
		if _, err := res.Write([]byte("Failed to parse body. Invalid book format")); err != nil {
			fmt.Println(err)
		}
		return
	}

	newBook.UUID = uuid
	done, err := db.BookDB.UpdateBook(&newBook)

	if err != nil || !done {
		res.WriteHeader(http.StatusNotFound)
		if _, err := res.Write([]byte("Book not found")); err != nil {
			fmt.Println(err)
		}
		return
	}

	res.WriteHeader(http.StatusCreated)
	if _, err := res.Write([]byte("Successfully updated the book with UUID: " + uuid)); err != nil {
		fmt.Println(err)
	}
	prometheusMetrics.BookUpdateCounter.Inc()
}
