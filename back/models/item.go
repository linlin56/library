package models

// IItem is the product interface of the abstract factory.
// Book and Journal satisfy it automatically by embedding Item.
type IItem interface {
	GetId() int
	GetTitle() string
	GetDeposit() float64
}

type Item struct {
	Id      int     `json:"id"`
	Deposit float64 `json:"deposit"`
	Title   string  `json:"title"`
	Copies  int     `json:"copies,omitempty"`
}

func (i *Item) GetId() int          { return i.Id }
func (i *Item) GetTitle() string    { return i.Title }
func (i *Item) GetDeposit() float64 { return i.Deposit }

func NewItem(id int, deposit float64, title string, copies int) *Item {
	return &Item{
		Id:      id,
		Deposit: deposit,
		Title:   title,
		Copies:  copies,
	}
}
