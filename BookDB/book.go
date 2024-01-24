package BookDB

import (
	"encoding/json"
	"errors"
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
	sync.Mutex
	Books map[string]*Book
}

var BookDB bookDBType

func Init() {
	BookDB = bookDBType{Books: make(map[string]*Book)}
}

func (this *bookDBType) UUID_Exists(uuid *string) bool {
	(*this).Lock()
	_, ok := (*this).Books[*uuid]
	(*this).Unlock()
	return ok
}

// returns "UUID", error
func (this *bookDBType) AddBook(newBook *Book) (string, error) {
	if newBook == nil {
		return "", errors.New("Book is Null")
	}

	var newUUID string
	for newUUID = (uuid.New()).String(); (*this).UUID_Exists(&newUUID); {
	}
	newBook.UUID = newUUID

	(*this).Lock()
	(*this).Books[newUUID] = newBook
	(*this).Unlock()
	return newUUID, nil
}

func (this *bookDBType) DeleteBook(UUID *string) (bool, error) {
	if !this.UUID_Exists(UUID) {
		return false, errors.New("No books found with the given UUID")
	}

	(*this).Lock()
	delete((*this).Books, *UUID)
	(*this).Unlock()

	return true, nil
}

func (this *bookDBType) UpdateBook(updatedBook *Book) (bool, error) {
	if !this.UUID_Exists(&(*updatedBook).UUID) {
		return false, errors.New("No books found with the given UUID. Please add the book before updating")
	}

	(*this).Lock()
	(*this).Books[(*updatedBook).UUID] = updatedBook
	(*this).Unlock()

	return true, nil
}

func (this *bookDBType) getAllBooks() ([][]byte, error) {
	(*this).Lock()
	defer (*this).Unlock()

	var bookList [][]byte

	for i, _ := range (*this).Books {
		tmp, err := json.MarshalIndent((*this).Books[i], " ", "   ")
		if err != nil {
			return [][]byte{}, err
		}
		bookList = append(bookList, tmp)
	}

	return bookList, nil
}

// /
func (this *bookDBType) GetBookList() string {
	bytes, err := (*this).getAllBooks()
	if err != nil {
		return "[]"
	}

	stringList := []string{}

	for _, v := range bytes {
		stringList = append(stringList, string(v))
	}

	return string("[\n" + strings.Join(stringList[:], ",\n") + "\n]")
}
