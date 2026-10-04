package models

import "time"

type User struct {
	Id             int       `json:"id"`
	LastName       string    `json:"last_name"`
	FirstName      string    `json:"first_name"`
	PhoneNumber    string    `json:"phone_number"`
	DepositBalance float64   `json:"deposit_balance"`
	TotalDeposit   float64   `json:"total_deposit"`
	Login          string    `json:"login"`
	Password       string    `json:"password"`
	BirthDate      time.Time `json:"birth_date"`
	Email          string    `json:"email"`
	AddressId      *int      `json:"address_id,omitempty"`
	Loans          []*Copy   `json:"loans,omitempty"`
}

func NewUser(id int, lastName string, firstName string, phoneNumber string, depositBalance float64, login string, password string, birthDate time.Time, email string) *User {
	return &User{
		Id:             id,
		LastName:       lastName,
		FirstName:      firstName,
		PhoneNumber:    phoneNumber,
		DepositBalance: depositBalance,
		Login:          login,
		Password:       password,
		BirthDate:      birthDate,
		Email:          email,
	}
}
