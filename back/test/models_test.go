package test

import (
	"library/models"
	"testing"
	"time"
)

// These tests check that the constructors behave as expected
// ===================== Author =====================

func TestNewAuthor_FieldsInitialized(t *testing.T) {
	a := models.NewAuthor(1, "Hugo", "Victor")
	if a.Id != 1 {
		t.Errorf("Id: expected 1, got %d", a.Id)
	}
	if a.LastName != "Hugo" {
		t.Errorf("LastName: expected Hugo, got %s", a.LastName)
	}
	if a.FirstName != "Victor" {
		t.Errorf("FirstName: expected Victor, got %s", a.FirstName)
	}
}

func TestNewAuthor_EmptyName(t *testing.T) {
	a := models.NewAuthor(0, "", "")
	if a.LastName != "" {
		t.Errorf("Empty LastName expected, got %s", a.LastName)
	}
}

// ===================== Item =====================

func TestNewItem_FieldsInitialized(t *testing.T) {
	i := models.NewItem(10, 3.50, "Les Miserables", 5)
	if i.Id != 10 {
		t.Errorf("Id: expected 10, got %d", i.Id)
	}
	if i.Deposit != 3.50 {
		t.Errorf("Deposit: expected 3.50, got %.2f", i.Deposit)
	}
	if i.Title != "Les Miserables" {
		t.Errorf("Wrong title: %s", i.Title)
	}
	if i.Copies != 5 {
		t.Errorf("Copies: expected 5, got %d", i.Copies)
	}
}

func TestItem_GetId(t *testing.T) {
	i := models.NewItem(42, 0, "", 0)
	if i.GetId() != 42 {
		t.Errorf("GetId: expected 42, got %d", i.GetId())
	}
}

func TestItem_GetTitle(t *testing.T) {
	i := models.NewItem(0, 0, "Test Title", 0)
	if i.GetTitle() != "Test Title" {
		t.Errorf("GetTitle: expected Test Title, got %s", i.GetTitle())
	}
}

func TestItem_GetDeposit(t *testing.T) {
	i := models.NewItem(0, 7.25, "", 0)
	if i.GetDeposit() != 7.25 {
		t.Errorf("GetDeposit: expected 7.25, got %.2f", i.GetDeposit())
	}
}

func TestItem_SatisfiesIItem(t *testing.T) {
	i := models.NewItem(1, 2.0, "Test", 1)
	var _ models.IItem = i
}

// ===================== Book =====================

func TestNewBook_FieldsInitialized(t *testing.T) {
	author := *models.NewAuthor(1, "Tolkien", "J.R.R.")
	b := models.NewBook(3, 4.0, "The Lord of the Rings", 2, author, "978-0618260584")
	if b.Id != 3 {
		t.Errorf("Id: expected 3, got %d", b.Id)
	}
	if b.Deposit != 4.0 {
		t.Errorf("Deposit: expected 4.0, got %.2f", b.Deposit)
	}
	if b.Title != "The Lord of the Rings" {
		t.Errorf("Wrong title: %s", b.Title)
	}
	if b.Copies != 2 {
		t.Errorf("Copies: expected 2, got %d", b.Copies)
	}
	if b.Isbn != "978-0618260584" {
		t.Errorf("Wrong Isbn: %s", b.Isbn)
	}
}

func TestNewBook_AuthorNotNil(t *testing.T) {
	author := *models.NewAuthor(5, "Zola", "Emile")
	b := models.NewBook(1, 2.0, "Germinal", 3, author, "978-2070413430")
	if b.Author == nil {
		t.Fatal("Author must not be nil")
	}
	if b.Author.LastName != "Zola" {
		t.Errorf("Author.LastName: expected Zola, got %s", b.Author.LastName)
	}
}

func TestBook_SatisfiesIItem(t *testing.T) {
	author := *models.NewAuthor(1, "Hugo", "Victor")
	b := models.NewBook(1, 2.0, "Notre-Dame de Paris", 1, author, "978-0140443530")
	var _ models.IItem = b
}

func TestBook_GetTitleInheritedFromItem(t *testing.T) {
	author := *models.NewAuthor(1, "Camus", "Albert")
	b := models.NewBook(1, 1.0, "L Etranger", 1, author, "978-2070360024")
	if b.GetTitle() != "L Etranger" {
		t.Errorf("Wrong GetTitle: %s", b.GetTitle())
	}
}

// ===================== Journal =====================

func TestNewJournal_FieldsInitialized(t *testing.T) {
	date, _ := time.Parse("2006-01-02", "2024-03-01")
	j := models.NewJournal(7, 1.50, "Science et Vie", 3, 42, date)
	if j.Id != 7 {
		t.Errorf("Id: expected 7, got %d", j.Id)
	}
	if j.Deposit != 1.50 {
		t.Errorf("Deposit: expected 1.50, got %.2f", j.Deposit)
	}
	if j.Title != "Science et Vie" {
		t.Errorf("Wrong title: %s", j.Title)
	}
	if j.IssueNumber != 42 {
		t.Errorf("IssueNumber: expected 42, got %d", j.IssueNumber)
	}
	if !j.PublicationDate.Equal(date) {
		t.Errorf("Wrong PublicationDate: expected %v, got %v", date, j.PublicationDate)
	}
}

func TestJournal_SatisfiesIItem(t *testing.T) {
	j := models.NewJournal(1, 1.0, "Test Journal", 0, 1, time.Now())
	var _ models.IItem = j
}

func TestJournal_GetDeposit(t *testing.T) {
	j := models.NewJournal(1, 2.75, "Test", 0, 1, time.Now())
	if j.GetDeposit() != 2.75 {
		t.Errorf("GetDeposit: expected 2.75, got %.2f", j.GetDeposit())
	}
}

// ===================== User =====================

func TestNewUser_FieldsInitialized(t *testing.T) {
	birthDate, _ := time.Parse("2006-01-02", "1990-04-15")
	u := models.NewUser(1, "Martin", "Jean", "0612345678", 20.0, "jmartin", "password", birthDate, "jean@example.com")
	if u.Id != 1 {
		t.Errorf("Id: expected 1, got %d", u.Id)
	}
	if u.LastName != "Martin" {
		t.Errorf("LastName: expected Martin, got %s", u.LastName)
	}
	if u.DepositBalance != 20.0 {
		t.Errorf("DepositBalance: expected 20.0, got %.2f", u.DepositBalance)
	}
	if u.Login != "jmartin" {
		t.Errorf("Login: expected jmartin, got %s", u.Login)
	}
	if u.Email != "jean@example.com" {
		t.Errorf("Wrong Email: %s", u.Email)
	}
}

func TestNewUser_AddressIdNilByDefault(t *testing.T) {
	u := models.NewUser(0, "", "", "", 0, "", "", time.Time{}, "")
	if u.AddressId != nil {
		t.Error("AddressId must be nil by default")
	}
}

// ===================== Student =====================

func TestNewStudent_FieldsInitialized(t *testing.T) {
	s := models.NewStudent("Dupont", "Marie", "mdupont", "password", "L3")
	if s.LastName != "Dupont" {
		t.Errorf("LastName: expected Dupont, got %s", s.LastName)
	}
	if s.FirstName != "Marie" {
		t.Errorf("FirstName: expected Marie, got %s", s.FirstName)
	}
	if s.Login != "mdupont" {
		t.Errorf("Login: expected mdupont, got %s", s.Login)
	}
	if s.StudyYear != "L3" {
		t.Errorf("StudyYear: expected L3, got %s", s.StudyYear)
	}
}

func TestStudent_InheritsUser(t *testing.T) {
	s := models.NewStudent("Test", "User", "tuser", "password", "M1")
	s.DepositBalance = 15.0
	if s.DepositBalance != 15.0 {
		t.Errorf("Inherited DepositBalance: expected 15.0, got %.2f", s.DepositBalance)
	}
}

// ===================== Teacher =====================

func TestNewTeacher_FieldsInitialized(t *testing.T) {
	dept := models.NewDepartment(1, "Informatique")
	teacher := models.NewTeacher("Roche", "Pierre", "proche", "password", dept)
	if teacher.LastName != "Roche" {
		t.Errorf("LastName: expected Roche, got %s", teacher.LastName)
	}
	if teacher.FirstName != "Pierre" {
		t.Errorf("FirstName: expected Pierre, got %s", teacher.FirstName)
	}
	if teacher.Login != "proche" {
		t.Errorf("Login: expected proche, got %s", teacher.Login)
	}
}

func TestNewTeacher_DepartmentNotNil(t *testing.T) {
	dept := models.NewDepartment(2, "Physique")
	teacher := models.NewTeacher("Bernard", "Claire", "cbernard", "password", dept)
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

func TestNewTeacher_WithoutDepartment(t *testing.T) {
	teacher := models.NewTeacher("Test", "User", "tuser", "password", nil)
	if teacher.Department != nil {
		t.Error("Department must be nil when nil is passed")
	}
}

// ===================== Librarian =====================

func TestNewLibrarian_FieldsInitialized(t *testing.T) {
	l := models.NewLibrarian(1, "Petit", "Valerie", "vpetit", "password")
	if l.Id != 1 {
		t.Errorf("Id: expected 1, got %d", l.Id)
	}
	if l.LastName != "Petit" {
		t.Errorf("LastName: expected Petit, got %s", l.LastName)
	}
	if l.Login != "vpetit" {
		t.Errorf("Login: expected vpetit, got %s", l.Login)
	}
	if l.Password != "password" {
		t.Errorf("Password: expected password, got %s", l.Password)
	}
}

// ===================== Department =====================

func TestNewDepartment_FieldsInitialized(t *testing.T) {
	d := models.NewDepartment(3, "Mathematiques")
	if d.Id != 3 {
		t.Errorf("Id: expected 3, got %d", d.Id)
	}
	if d.Name != "Mathematiques" {
		t.Errorf("Name: expected Mathematiques, got %s", d.Name)
	}
}

func TestNewDepartment_ZeroIdIsValid(t *testing.T) {
	d := models.NewDepartment(0, "No ID")
	if d.Id != 0 {
		t.Errorf("Id: expected 0, got %d", d.Id)
	}
}

// ===================== Copy =====================

func TestNewCopy_FieldsInitialized(t *testing.T) {
	c := models.NewCopy("EX-0001", 15)
	if c.Barcode != "EX-0001" {
		t.Errorf("Barcode: expected EX-0001, got %s", c.Barcode)
	}
	if c.LoanPeriodDays != 15 {
		t.Errorf("LoanPeriodDays: expected 15, got %d", c.LoanPeriodDays)
	}
}

func TestNewCopy_NotBorrowedByDefault(t *testing.T) {
	c := models.NewCopy("EX-0042", 7)
	if c.IsBorrowed {
		t.Error("IsBorrowed must be false by default")
	}
}

func TestNewCopy_EmptyDates(t *testing.T) {
	c := models.NewCopy("EX-0010", 15)
	if !c.LoanStartDate.IsZero() {
		t.Error("LoanStartDate must be zero by default")
	}
	if !c.LoanDueDate.IsZero() {
		t.Error("LoanDueDate must be zero by default")
	}
}
