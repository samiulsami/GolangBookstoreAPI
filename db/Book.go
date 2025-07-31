package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/google/uuid"
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

var (
	BookDB          bookDBType
	ErrBookNotFound = errors.New("no books found with the given UUID")
)

func init() {
	BookDB = bookDBType{books: make(map[string]*Book)}
}

func (book *Book) IsValid() bool {
	if book == nil {
		return false
	}

	if book.Name == "" || book.ISBN == "" || book.AuthorList == nil || book.PublishDate == "" {
		return false
	}

	return true
}

func (bookdb *bookDBType) uuidExists(uuid *string) bool {
	_, ok := bookdb.books[*uuid]
	return ok
}

func (bookdb *bookDBType) BookExists(uuid string) bool {
	bookdb.RLock()
	defer bookdb.RUnlock()
	return bookdb.uuidExists(&uuid)
}

// returns "UUID"
func (bookdb *bookDBType) AddBook(newBook *Book) string {
	bookdb.Lock()
	defer bookdb.Unlock()

	newUUID := (uuid.New()).String()
	for ; bookdb.uuidExists(&newUUID); newUUID = (uuid.New()).String() {
	}
	newBook.UUID = newUUID
	bookdb.books[newUUID] = newBook

	return newUUID
}

func (bookdb *bookDBType) DeleteBook(uuid string) (bool, error) {
	bookdb.Lock()
	defer bookdb.Unlock()

	if !bookdb.uuidExists(&uuid) {
		return false, ErrBookNotFound
	}

	delete(bookdb.books, uuid)
	return true, nil
}

func (bookdb *bookDBType) UpdateBook(updatedBook *Book) (bool, error) {
	bookdb.Lock()
	defer bookdb.Unlock()

	if !bookdb.uuidExists(&(*updatedBook).UUID) {
		return false, ErrBookNotFound
	}

	bookdb.books[(*updatedBook).UUID] = updatedBook
	return true, nil
}

// /returns a single book as json object, if it exists
func (bookdb *bookDBType) GetBook(uuid string) ([]byte, error) {
	bookdb.RLock()
	defer bookdb.RUnlock()

	if !(bookdb.uuidExists(&uuid)) {
		return []byte(""), ErrBookNotFound
	}

	tmp, err := json.MarshalIndent((*bookdb).books[uuid], " ", "   ")
	if err != nil {
		return []byte(""), err
	}

	return tmp, nil
}

func (bookdb *bookDBType) getAllBooks() ([][]byte, error) {
	bookdb.RLock()
	defer bookdb.RUnlock()

	var bookList [][]byte

	for i := range bookdb.books {
		tmp, err := json.MarshalIndent(bookdb.books[i], " ", "   ")
		if err != nil {
			return [][]byte{}, err
		}
		bookList = append(bookList, tmp)
	}

	return bookList, nil
}

// returns a json array of books as a byte slice
func (bookdb *bookDBType) GetBookList() []byte {
	bytes, err := bookdb.getAllBooks()
	if err != nil {
		fmt.Println(err)
		return []byte("[]")
	}

	stringList := []string{}

	for _, v := range bytes {
		stringList = append(stringList, string(v))
	}

	return []byte("[\n" + strings.Join(stringList[:], ",\n") + "\n]")
}
