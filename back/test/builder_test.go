package test

import (
	"library/builder"
	"library/models"
	"testing"
	"time"
)

// Tests for the builder design pattern
// ===================== UserBuilder =====================

func TestUserBuilder_Build_SimpleFields(t *testing.T) {
	birthDate := time.Date(1990, 5, 12, 0, 0, 0, 0, time.UTC)
	u := builder.NewUserBuilder().
		WithLastName("Martin").
		WithFirstName("Jean").
		WithLogin("jmartin").
		WithPassword("password").
		WithEmail("jean@example.com").
		WithPhoneNumber("0612345678").
		WithDepositBalance(20.0).
		WithBirthDate(birthDate).
		Build()
	if u.LastName != "Martin" {
		t.Errorf("LastName: expected Martin, got %s", u.LastName)
	}
	if u.FirstName != "Jean" {
		t.Errorf("FirstName: expected Jean, got %s", u.FirstName)
	}
	if u.Login != "jmartin" {
		t.Errorf("Login: expected jmartin, got %s", u.Login)
	}
	if u.Password != "password" {
		t.Errorf("Wrong Password: %s", u.Password)
	}
	if u.Email != "jean@example.com" {
		t.Errorf("Wrong Email: %s", u.Email)
	}
	if u.PhoneNumber != "0612345678" {
		t.Errorf("Wrong PhoneNumber: %s", u.PhoneNumber)
	}
	if u.DepositBalance != 20.0 {
		t.Errorf("DepositBalance: expected 20.0, got %.2f", u.DepositBalance)
	}
	if !u.BirthDate.Equal(birthDate) {
		t.Errorf("Wrong BirthDate")
	}
}

func TestUserBuilder_Build_ReturnsUserPointer(t *testing.T) {
	u := builder.NewUserBuilder().WithLastName("Test").Build()
	_ = u
}

func TestUserBuilder_Build_EmptyDefaults(t *testing.T) {
	u := builder.NewUserBuilder().Build()
	if u.LastName != "" {
		t.Errorf("LastName should be empty by default, got %s", u.LastName)
	}
	if u.DepositBalance != 0 {
		t.Errorf("DepositBalance should be 0 by default, got %.2f", u.DepositBalance)
	}
}

func TestUserBuilder_Chaining_ReturnsSameBuilder(t *testing.T) {
	b := builder.NewUserBuilder()
	u := b.WithLastName("Test").WithFirstName("User").Build()
	if u.LastName != "Test" {
		t.Errorf("LastName after chaining: expected Test, got %s", u.LastName)
	}
}

// ===================== StudentBuilder =====================

func TestStudentBuilder_Build_AllFields(t *testing.T) {
	birthDate := time.Date(2000, 9, 1, 0, 0, 0, 0, time.UTC)
	s := builder.NewStudentBuilder().
		WithLastName("Dupont").
		WithFirstName("Marie").
		WithLogin("mdupont").
		WithPassword("password").
		WithEmail("marie@example.com").
		WithPhoneNumber("0698765432").
		WithDepositBalance(20.0).
		WithBirthDate(birthDate).
		WithStudyYear("L3").
		Build()
	if s.LastName != "Dupont" {
		t.Errorf("LastName: expected Dupont, got %s", s.LastName)
	}
	if s.StudyYear != "L3" {
		t.Errorf("StudyYear: expected L3, got %s", s.StudyYear)
	}
	if s.DepositBalance != 20.0 {
		t.Errorf("DepositBalance: expected 20.0, got %.2f", s.DepositBalance)
	}
}

func TestStudentBuilder_Build_ReturnsStudentPointer(t *testing.T) {
	s := builder.NewStudentBuilder().WithLastName("Test").Build()
	_ = s
}

func TestStudentBuilder_Build_StudyYearEmptyByDefault(t *testing.T) {
	s := builder.NewStudentBuilder().WithLastName("Test").Build()
	if s.StudyYear != "" {
		t.Errorf("StudyYear should be empty by default, got %s", s.StudyYear)
	}
}

func TestStudentBuilder_InheritsUserFields(t *testing.T) {
	s := builder.NewStudentBuilder().
		WithLastName("Bernard").
		WithFirstName("Alice").
		WithDepositBalance(15.0).
		Build()
	if s.LastName != "Bernard" {
		t.Errorf("LastName: expected Bernard, got %s", s.LastName)
	}
	if s.DepositBalance != 15.0 {
		t.Errorf("User.DepositBalance: expected 15.0, got %.2f", s.DepositBalance)
	}
}

// ===================== TeacherBuilder =====================

func TestTeacherBuilder_Build_WithoutDepartment(t *testing.T) {
	teacher := builder.NewTeacherBuilder().
		WithLastName("Roche").
		WithFirstName("Pierre").
		WithLogin("proche").
		WithPassword("password").
		Build()
	if teacher.LastName != "Roche" {
		t.Errorf("LastName: expected Roche, got %s", teacher.LastName)
	}
	if teacher.Department != nil {
		t.Error("Department must be nil when not set")
	}
}

func TestTeacherBuilder_Build_WithDepartment(t *testing.T) {
	dept := models.NewDepartment(2, "Physique")
	teacher := builder.NewTeacherBuilder().
		WithLastName("Bernard").
		WithFirstName("Claire").
		WithLogin("cbernard").
		WithDepartment(dept).
		Build()
	if teacher.Department == nil {
		t.Fatal("Department must not be nil")
	}
	if teacher.Department.Id != 2 {
		t.Errorf("Department.Id: expected 2, got %d", teacher.Department.Id)
	}
	if teacher.Department.Name != "Physique" {
		t.Errorf("Department.Name: expected Physique, got %s", teacher.Department.Name)
	}
}

func TestTeacherBuilder_Build_ReturnsTeacherPointer(t *testing.T) {
	teacher := builder.NewTeacherBuilder().WithLastName("Test").Build()
	_ = teacher
}

func TestTeacherBuilder_Build_UserFieldsPropagated(t *testing.T) {
	teacher := builder.NewTeacherBuilder().
		WithLastName("Durand").
		WithFirstName("Luc").
		WithEmail("luc@example.com").
		WithDepositBalance(20.0).
		Build()
	if teacher.LastName != "Durand" {
		t.Errorf("User.LastName: expected Durand, got %s", teacher.LastName)
	}
	if teacher.Email != "luc@example.com" {
		t.Errorf("Wrong User.Email: %s", teacher.Email)
	}
}

func TestTeacherBuilder_Build_BirthDate(t *testing.T) {
	birthDate := time.Date(1975, 3, 20, 0, 0, 0, 0, time.UTC)
	teacher := builder.NewTeacherBuilder().
		WithLastName("Lefebvre").
		WithBirthDate(birthDate).
		Build()
	if !teacher.BirthDate.Equal(birthDate) {
		t.Errorf("Wrong BirthDate")
	}
}
