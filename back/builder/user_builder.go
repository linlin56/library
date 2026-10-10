package builder

import (
	"library/models"
	"time"
)

// Builder pattern
type UserBuilder struct {
	u models.User
}

func NewUserBuilder() *UserBuilder {
	return &UserBuilder{}
}

func (b *UserBuilder) WithLastName(lastName string) *UserBuilder {
	b.u.LastName = lastName
	return b
}

func (b *UserBuilder) WithFirstName(firstName string) *UserBuilder {
	b.u.FirstName = firstName
	return b
}

func (b *UserBuilder) WithLogin(login string) *UserBuilder {
	b.u.Login = login
	return b
}

func (b *UserBuilder) WithPassword(password string) *UserBuilder {
	b.u.Password = password
	return b
}

func (b *UserBuilder) WithEmail(email string) *UserBuilder {
	b.u.Email = email
	return b
}

func (b *UserBuilder) WithPhoneNumber(phone string) *UserBuilder {
	b.u.PhoneNumber = phone
	return b
}

func (b *UserBuilder) WithDepositBalance(balance float64) *UserBuilder {
	b.u.DepositBalance = balance
	return b
}

func (b *UserBuilder) WithBirthDate(date time.Time) *UserBuilder {
	b.u.BirthDate = date
	return b
}

func (b *UserBuilder) WithAddressId(id int) *UserBuilder {
	b.u.AddressId = &id
	return b
}

func (b *UserBuilder) Build() *models.User {
	return &b.u
}
