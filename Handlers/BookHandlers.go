package Handlers

import (
	"GoBookstoreAPI/DB"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func AddBook(res http.ResponseWriter, req *http.Request) {
	var newBook DB.Book
	err := json.NewDecoder(req.Body).Decode(&newBook)

	if err != nil || !DB.ValidBook(&newBook) {
		res.Write([]byte("Failed to parse body. Invalid book format"))
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	uuid := DB.BookDB.AddBook(&newBook)
	res.Write([]byte("Book added. UUID: " + uuid))
	res.WriteHeader(http.StatusCreated)
}

func GetBook(res http.ResponseWriter, req *http.Request) {
	uuid := chi.URLParam(req, "id")
	body, err := DB.BookDB.GetBook(uuid)

	if err != nil {
		fmt.Println(err)
		res.WriteHeader(http.StatusNotFound)
		return
	}

	res.Write(body)
	res.WriteHeader(http.StatusOK)
}

func GetAllBooks(res http.ResponseWriter, req *http.Request) {
	res.Write(DB.BookDB.GetBookList())
	res.WriteHeader(http.StatusOK)
}

func DeleteBook(res http.ResponseWriter, req *http.Request) {
	uuid := chi.URLParam(req, "id")
	done, err := DB.BookDB.DeleteBook(uuid)

	if err != nil || !done {
		fmt.Println(err)
		res.WriteHeader(http.StatusForbidden)
		return
	}

	res.Write([]byte("Successfully deleted the book with UUID: " + uuid))
	res.WriteHeader(http.StatusOK)
}

func UpdateBook(res http.ResponseWriter, req *http.Request) {
	uuid := chi.URLParam(req, "id")
	var newBook DB.Book
	err := json.NewDecoder(req.Body).Decode(&newBook)

	if err != nil || !DB.ValidBook(&newBook) {
		res.Write([]byte("Failed to parse body. Invalid book format"))
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	newBook.UUID = uuid
	done, err := DB.BookDB.UpdateBook(&newBook)

	if err != nil || !done {
		res.Write([]byte("Book not found"))
		res.WriteHeader(http.StatusForbidden)
		return
	}

	res.Write([]byte("Successfully updated the book with UUID: " + uuid))
	res.WriteHeader(http.StatusCreated)
}
