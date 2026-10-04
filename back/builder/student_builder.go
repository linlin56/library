package builder

import (
	"library/models"
	"time"
)

// Builder pattern
// StudentBuilder extends UserBuilder to build a Student.
type StudentBuilder struct {
	UserBuilder
	studyYear string
}

func NewStudentBuilder() *StudentBuilder {
	return &StudentBuilder{}
}

// The UserBuilder With* methods return *UserBuilder,
// so each one is redefined to keep the *StudentBuilder type.

func (b *StudentBuilder) WithLastName(lastName string) *StudentBuilder {
	b.u.LastName = lastName
	return b
}

func (b *StudentBuilder) WithFirstName(firstName string) *StudentBuilder {
	b.u.FirstName = firstName
	return b
}

func (b *StudentBuilder) WithLogin(login string) *StudentBuilder {
	b.u.Login = login
	return b
}

func (b *StudentBuilder) WithPassword(password string) *StudentBuilder {
	b.u.Password = password
	return b
}

func (b *StudentBuilder) WithEmail(email string) *StudentBuilder {
	b.u.Email = email
	return b
}

func (b *StudentBuilder) WithPhoneNumber(phone string) *StudentBuilder {
	b.u.PhoneNumber = phone
	return b
}

func (b *StudentBuilder) WithDepositBalance(balance float64) *StudentBuilder {
	b.u.DepositBalance = balance
	return b
}

func (b *StudentBuilder) WithBirthDate(date time.Time) *StudentBuilder {
	b.u.BirthDate = date
	return b
}

func (b *StudentBuilder) WithStudyYear(year string) *StudentBuilder {
	b.studyYear = year
	return b
}

func (b *StudentBuilder) Build() *models.Student {
	return &models.Student{
		User:      b.u,
		StudyYear: b.studyYear,
	}
}
