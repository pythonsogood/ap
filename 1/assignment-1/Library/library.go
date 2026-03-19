package Library

import (
	"errors"
	"fmt"
)

type Book struct {
	Id         string
	Title      string
	Author     string
	IsBorrowed bool
}

type Library struct {
	Books map[string]Book
}

func (l *Library) AddBook(id string, title string, author string) (*Book, error) {
	book, ok := l.Books[title]

	if ok {
		return nil, errors.New("Book already exists")
	}

	book = Book{Id: id, Title: title, Author: author, IsBorrowed: false}

	l.Books[book.Id] = book

	return &book, nil
}

func (l *Library) BorrowBook(id string) (*Book, error) {
	book, ok := l.Books[id]

	if !ok {
		return nil, errors.New("Book not found")
	}

	if book.IsBorrowed {
		return nil, errors.New("Book is already borrowed")
	}

	book.IsBorrowed = true

	l.Books[book.Id] = book

	return &book, nil
}

func (l *Library) ReturnBook(id string) error {
	book, ok := l.Books[id]

	if !ok {
		return errors.New("Book not found")
	}

	if !book.IsBorrowed {
		return errors.New("Book is not borrowed")
	}

	book.IsBorrowed = false

	l.Books[book.Id] = book

	return nil
}

func (l *Library) ListAvailableBooks() []Book {
	var books []Book

	for _, book := range l.Books {
		if !book.IsBorrowed {
			books = append(books, book)
		}
	}

	return books
}

func NewLibrary() *Library {
	return &Library{Books: make(map[string]Book)}
}

func CLI() {
	library := NewLibrary()

	for {
		fmt.Println("\n--- Library ---")
		fmt.Println("[1] Add Book")
		fmt.Println("[2] Borrow Book")
		fmt.Println("[3] Return Book")
		fmt.Println("[4] List Available Books")
		fmt.Println("[0] Return")
		fmt.Print("\n>>> ")

		var choice int

		fmt.Scanln(&choice)

		switch choice {
		case 1:
			var id, title, author string

			fmt.Println("Enter Book ID:")
			fmt.Scanln(&id)

			fmt.Println("Enter Book Title:")
			fmt.Scanln(&title)

			fmt.Println("Enter Book Author:")
			fmt.Scanln(&author)

			_, err := library.AddBook(id, title, author)

			if err != nil {
				fmt.Println(err)
			}

		case 2:
			var id string

			fmt.Println("Enter Book ID:")
			fmt.Scanln(&id)

			_, err := library.BorrowBook(id)

			if err != nil {
				fmt.Println(err)
			}

		case 3:
			var id string

			fmt.Println("Enter Book ID:")
			fmt.Scanln(&id)

			err := library.ReturnBook(id)

			if err != nil {
				fmt.Println(err)
			}

		case 4:
			books := library.ListAvailableBooks()

			if len(books) == 0 {
				fmt.Println("No available books")
				continue
			}

			for _, book := range books {
				fmt.Printf("ID: %s, Title: %s, Author: %s\n", book.Id, book.Title, book.Author)
			}

		case 0:
			return

		default:
			fmt.Println("Invalid choice")
		}
	}
}
