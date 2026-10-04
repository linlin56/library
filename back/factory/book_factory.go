package factory

import "library/models"

// BookFactory is the concrete factory for books.
type BookFactory struct {
	Title   string
	Deposit float64
	Isbn    string
	Author  *models.Author
}

// CreateItem implements ItemFactory and returns a *models.Book.
func (f *BookFactory) CreateItem() models.IItem {
	return models.NewBook(0, f.Deposit, f.Title, 0, *f.Author, f.Isbn)
}
