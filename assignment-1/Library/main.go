package Library

type Book struct {
	ID         int
	Title      string
	Author     string
	IsBorrowed bool
}

type Library struct {
	Books map[string]Book
}
