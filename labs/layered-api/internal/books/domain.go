package books

// Book is the domain representation owned by the application, not by HTTP.
type Book struct {
	ID     int64
	Title  string
	Author string
}

// CreateBookInput contains the business data required to create a book.
type CreateBookInput struct {
	Title  string
	Author string
}
