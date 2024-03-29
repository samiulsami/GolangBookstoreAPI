package db

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
var ErrBookNotFound = errors.New("No books found with the given UUID")

func init() {
	BookDB = bookDBType{books: make(map[string]*Book)}
}

func (this *Book) IsValid() bool {
	if this == nil {
		return false
	}

	if this.Name == "" || this.ISBN == "" || this.AuthorList == nil || this.PublishDate == "" {
		return false
	}

	return true
}

func (this *bookDBType) uuidExists(uuid *string) bool {
	_, ok := this.books[*uuid]
	return ok
}

func (this *bookDBType) BookExists(uuid string) bool {
	this.RLock()
	defer this.RUnlock()
	return this.uuidExists(&uuid)
}

// returns "UUID"
func (this *bookDBType) AddBook(newBook *Book) string {
	this.Lock()
	defer this.Unlock()

	var newUUID string = (uuid.New()).String()
	for ; this.uuidExists(&newUUID); newUUID = (uuid.New()).String() {
	}
	newBook.UUID = newUUID
	this.books[newUUID] = newBook

	return newUUID
}

func (this *bookDBType) DeleteBook(uuid string) (bool, error) {
	this.Lock()
	defer this.Unlock()

	if !this.uuidExists(&uuid) {
		return false, ErrBookNotFound
	}

	delete(this.books, uuid)
	return true, nil
}

func (this *bookDBType) UpdateBook(updatedBook *Book) (bool, error) {
	this.Lock()
	defer this.Unlock()

	if !this.uuidExists(&(*updatedBook).UUID) {
		return false, ErrBookNotFound
	}

	this.books[(*updatedBook).UUID] = updatedBook
	return true, nil
}

// /returns a single book as json object, if it exists
func (this *bookDBType) GetBook(uuid string) ([]byte, error) {
	this.RLock()
	defer this.RUnlock()

	if !(this.uuidExists(&uuid)) {
		return []byte(""), ErrBookNotFound
	}

	tmp, err := json.MarshalIndent((*this).books[uuid], " ", "   ")
	if err != nil {
		return []byte(""), err
	}

	return tmp, nil
}

func (this *bookDBType) getAllBooks() ([][]byte, error) {
	this.RLock()
	defer this.RUnlock()

	var bookList [][]byte

	for i, _ := range this.books {
		tmp, err := json.MarshalIndent(this.books[i], " ", "   ")
		if err != nil {
			return [][]byte{}, err
		}
		bookList = append(bookList, tmp)
	}

	return bookList, nil
}

// returns a json array of books as a byte slice
func (this *bookDBType) GetBookList() []byte {
	bytes, err := this.getAllBooks()
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
