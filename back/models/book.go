package models

type Book struct {
	Item
	AuthorId int     `json:"author_id,omitempty"`
	Author   *Author `json:"author,omitempty"`
	Isbn     string  `json:"isbn,omitempty"`
}

func NewBook(id int, deposit float64, title string, copies int, author Author, isbn string) *Book {
	return &Book{
		Item: Item{
			Id:      id,
			Deposit: deposit,
			Title:   title,
			Copies:  copies,
		},
		Author: &author,
		Isbn:   isbn,
	}
}
