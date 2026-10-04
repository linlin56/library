package factory

import (
	"time"

	"library/models"
)

// JournalFactory is the concrete factory for journals.
type JournalFactory struct {
	Title           string
	Deposit         float64
	IssueNumber     int
	PublicationDate time.Time
}

// CreateItem implements ItemFactory and returns a *models.Journal.
func (f *JournalFactory) CreateItem() models.IItem {
	return models.NewJournal(0, f.Deposit, f.Title, 0, f.IssueNumber, f.PublicationDate)
}
