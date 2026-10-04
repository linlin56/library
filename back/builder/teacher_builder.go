package builder

import (
	"library/models"
	"time"
)

// Builder pattern
// TeacherBuilder extends UserBuilder to build a Teacher.
type TeacherBuilder struct {
	UserBuilder
	department *models.Department
}

func NewTeacherBuilder() *TeacherBuilder {
	return &TeacherBuilder{}
}

func (b *TeacherBuilder) WithLastName(lastName string) *TeacherBuilder {
	b.u.LastName = lastName
	return b
}

func (b *TeacherBuilder) WithFirstName(firstName string) *TeacherBuilder {
	b.u.FirstName = firstName
	return b
}

func (b *TeacherBuilder) WithLogin(login string) *TeacherBuilder {
	b.u.Login = login
	return b
}

func (b *TeacherBuilder) WithPassword(password string) *TeacherBuilder {
	b.u.Password = password
	return b
}

func (b *TeacherBuilder) WithEmail(email string) *TeacherBuilder {
	b.u.Email = email
	return b
}

func (b *TeacherBuilder) WithPhoneNumber(phone string) *TeacherBuilder {
	b.u.PhoneNumber = phone
	return b
}

func (b *TeacherBuilder) WithDepositBalance(balance float64) *TeacherBuilder {
	b.u.DepositBalance = balance
	return b
}

func (b *TeacherBuilder) WithBirthDate(date time.Time) *TeacherBuilder {
	b.u.BirthDate = date
	return b
}

func (b *TeacherBuilder) WithDepartment(dept *models.Department) *TeacherBuilder {
	b.department = dept
	return b
}

func (b *TeacherBuilder) Build() *models.Teacher {
	return &models.Teacher{
		User:       b.u,
		Department: b.department,
	}
}
