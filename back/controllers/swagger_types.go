package controllers

// ErrorResponse describes a standard API error.
type ErrorResponse struct {
	Error string `json:"error"`
}

// MessageResponse describes a simple response with a message.
type MessageResponse struct {
	Message string `json:"message"`
}

// IDResponse describes a response that returns an identifier.
type IDResponse struct {
	ID int `json:"id"`
}

// CreateItemResponse describes the response when a book or a journal is created.
type CreateItemResponse struct {
	ID      int    `json:"id"`
	Message string `json:"message"`
}

// LoginConflictResponse describes the 409 response when a login already exists.
type LoginConflictResponse struct {
	Message       string `json:"message"`
	ExistingLogin string `json:"existing_login"`
}

// UpdateDepositRequest describes the deposit update payload.
type UpdateDepositRequest struct {
	TotalDeposit float64 `json:"total_deposit"`
}

// CreateCopyRequest describes the copy creation payload.
type CreateCopyRequest struct {
	Barcode        string `json:"barcode"`
	LoanPeriodDays int    `json:"loan_period_days"`
}
