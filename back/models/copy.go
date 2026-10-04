package models

import "time"

type Copy struct {
	Id             int       `json:"id"`
	IsBorrowed     bool      `json:"is_borrowed"`
	LoanStartDate  time.Time `json:"loan_start_date"`
	LoanDueDate    time.Time `json:"loan_due_date"`
	LoanPeriodDays int       `json:"loan_period_days"`
	Barcode        string    `json:"barcode"`
	BorrowerId     int       `json:"borrower_id"`
	Borrower       *User     `json:"borrower,omitempty"`
}

func NewCopy(barcode string, loanPeriodDays int) *Copy {
	return &Copy{
		IsBorrowed:     false,
		LoanPeriodDays: loanPeriodDays,
		Barcode:        barcode,
	}
}
