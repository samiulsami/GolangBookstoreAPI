package main

import (
	"GoBookstoreAPI/DB"
	"fmt"
)

func main() {
	fmt.Println("hehe")
	DB.Init()
	var books = DB.BookDB

	for i := 0; i < 5; i++ {
		myBook := DB.Book{
			Name: "hehe",
		}
		myBook.PublishDate = "abc" + fmt.Sprintf("%d%d%d%d", i, i, i, i) + "def"
		a, _ := books.AddBook(&myBook)
		fmt.Println(a)
		fmt.Println(books.Book_Exists(a))
	}

	fmt.Println(books.GetBookList())
}
