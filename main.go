package main

import (
	"GoBookstoreAPI/BookDB"
	"fmt"
)

func main() {
	fmt.Println("hehe")
	BookDB.Init()
	var books = BookDB.BookDB

	for i := 0; i < 5; i++ {
		myBook := BookDB.Book{
			Name: "hehe",
		}
		myBook.PublishDate = fmt.Sprintf("%d", i)
		a, _ := books.AddBook(&myBook)
		fmt.Println(a)
		fmt.Println(books.UUID_Exists(&a))
	}

	fmt.Println(len(books.Books))

	fmt.Println(books.GetAllBooksAsJson())

}
