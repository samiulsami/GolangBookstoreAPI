package DB

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"strings"
	"sync"
)

type Book struct {
	UUID        string   `json:"uuid"`
	Name        string   `json:"name"`
	AuthorList  []string `json:"authorList"`
	PublishDate string   `json:"publishDate"`
	ISBN        string   `json:"isbn"`
}

// {"UUID": book}
type bookDBType struct {
	sync.RWMutex
	books map[string]*Book
}

var BookDB bookDBType
var book_not_found_error string = "No books found with the given UUID."

func Init() {
	BookDB = bookDBType{books: make(map[string]*Book)}
}

func (this *bookDBType) uuid_Exists(uuid *string) bool {
	_, ok := (*this).books[*uuid]
	return ok
}

func (this *bookDBType) Book_Exists(uuid string) bool {
	(*this).RLock()
	defer (*this).RUnlock()
	return (*this).uuid_Exists(&uuid)
}

// returns "UUID", error
func (this *bookDBType) AddBook(newBook *Book) (string, error) {
	if newBook == nil {
		return "", errors.New("Book is Null")
	}

	(*this).Lock()
	defer (*this).Unlock()

	var newUUID string
	for newUUID = (uuid.New()).String(); (*this).uuid_Exists(&newUUID); {
	}
	newBook.UUID = newUUID
	(*this).books[newUUID] = newBook

	return newUUID, nil
}

func (this *bookDBType) DeleteBook(uuid string) (bool, error) {
	(*this).Lock()
	defer (*this).Unlock()

	if !this.uuid_Exists(&uuid) {
		return false, errors.New(book_not_found_error)
	}

	delete((*this).books, uuid)
	return true, nil
}

func (this *bookDBType) UpdateBook(updatedBook *Book) (bool, error) {
	(*this).Lock()
	defer (*this).Unlock()

	if !this.uuid_Exists(&(*updatedBook).UUID) {
		return false, errors.New(book_not_found_error)
	}

	(*this).books[(*updatedBook).UUID] = updatedBook
	return true, nil
}

// /returns a single book as json object, if it exists
func (this *bookDBType) GetBook(uuid string) (string, error) {
	(*this).RLock()
	defer (*this).RUnlock()

	if !(*this).uuid_Exists(&uuid) {
		return "", errors.New(book_not_found_error)
	}

	tmp, err := json.MarshalIndent((*this).books[uuid], " ", "   ")
	if err != nil {
		return "", err
	}

	return string(tmp), nil
}

func (this *bookDBType) getAllBooks() ([][]byte, error) {
	(*this).RLock()
	defer (*this).RUnlock()

	var bookList [][]byte

	for i, _ := range (*this).books {
		tmp, err := json.MarshalIndent((*this).books[i], " ", "   ")
		if err != nil {
			return [][]byte{}, err
		}
		bookList = append(bookList, tmp)
	}

	return bookList, nil
}

// returns a json array of books
func (this *bookDBType) GetBookList() string {
	bytes, err := (*this).getAllBooks()
	if err != nil {
		fmt.Println(err)
		return "[]"
	}

	stringList := []string{}

	for _, v := range bytes {
		stringList = append(stringList, string(v))
	}

	return string("[\n" + strings.Join(stringList[:], ",\n") + "\n]")
}
