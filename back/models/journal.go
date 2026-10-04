package models

import "time"

type Journal struct {
	Item
	IssueNumber     int       `json:"issue_number,omitempty"`
	PublicationDate time.Time `json:"publication_date,omitempty"`
}

func NewJournal(id int, deposit float64, title string, copies int, issueNumber int, publicationDate time.Time) *Journal {
	return &Journal{
		Item: Item{
			Id:      id,
			Deposit: deposit,
			Title:   title,
			Copies:  copies,
		},
		IssueNumber:     issueNumber,
		PublicationDate: publicationDate,
	}
}
