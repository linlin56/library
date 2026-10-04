// gen-uml command: generates the project's PlantUML sequence diagrams
// as well as a class diagram via goplantuml.
//
// Usage:
//
//	go run ./cmd/gen-uml            → writes to back/diagrams/
//	go run ./cmd/gen-uml -out /tmp  → custom directory
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"library/uml"
)

// This file generates sequence diagrams with a Go library,
// to compare them with my hand-made diagrams (Visual Paradigm) and check that everything matches.
func main() {
	out := flag.String("out", "diagrams", "Output directory for the .puml files")
	flag.Parse()

	if err := os.MkdirAll(*out, 0o755); err != nil {
		log.Fatalf("Unable to create the output directory: %v", err)
	}

	diagrams := map[string]uml.Diagram{
		// List of the diagrams to generate
		"01_login":       loginDiagram(),
		"02_loan":        loanDiagram(),
		"03_create_item": createItemDiagram(),
		// 04_search is maintained by hand — do not overwrite.
		// "04_search": searchDiagram(),
		"05_create_user":  createUserDiagram(),
		"06_list_loans":   listLoansDiagram(),
		"07_list_overdue": listOverdueDiagram(),
		"08_find_user":    findUserDiagram(),
		"09_use_cases":    useCaseDiagram(),
	}

	for name, d := range diagrams {
		// Write each diagram to a .puml file
		path := filepath.Join(*out, name+".puml")
		if err := d.WriteFile(path); err != nil {
			log.Fatalf("Error writing %s: %v", path, err)
		}
		fmt.Printf("OK  %s\n", path)
	}

	genClassDiagram(*out)
}

// genClassDiagram runs goplantuml to generate the class diagram.
func genClassDiagram(outDir string) {
	path, err := exec.LookPath("goplantuml")
	if err != nil {
		fmt.Println("INFO goplantuml not found in PATH — class diagram skipped.")
		fmt.Println("     Install it with: go install github.com/jfeliu007/goplantuml@latest")
		return
	}

	outFile := filepath.Join(outDir, "00_classes.puml")
	f, err := os.Create(outFile)
	if err != nil {
		log.Printf("Error creating %s: %v", outFile, err)
		return
	}
	defer func() { _ = f.Close() }()

	cmd := exec.Command(path, "-recursive", ".")
	cmd.Stdout = f
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Printf("goplantuml returned an error: %v", err)
		return
	}
	fmt.Printf("OK  %s\n", outFile)
}

// Diagram 1 — User login

func loginDiagram() *uml.SequenceDiagram {
	return uml.NewSequenceDiagram("User login").
		Actor("C", "Client").
		Participant("CTL", "LoginController").
		Participant("UR", "UserRepository").
		Participant("LR", "LibrarianRepository").
		Participant("DB", "PostgreSQL").
		Blank().
		Arrow("C", "CTL", "POST /api/login\\n{login, password}").
		Activate("CTL").
		Blank().
		Arrow("CTL", "UR", "FindByLogin(login)").
		Activate("UR").
		Arrow("UR", "DB", "SELECT * FROM ONLY users WHERE login = $1").
		Activate("DB").
		Blank().
		AltStart("User found").
		DashedArrow("DB", "UR", "user row").
		Deactivate("DB").
		DashedArrow("UR", "CTL", "*User, nil").
		Deactivate("UR").
		AltStart("Correct password").
		DashedArrow("CTL", "C", "200 {id, last_name, first_name, role: user}").
		Else("Wrong password").
		DashedArrow("CTL", "C", "401 Incorrect username or password").
		End().
		Blank().
		Else("User not found").
		DashedArrow("DB", "UR", "sql.ErrNoRows").
		Deactivate("DB").
		DashedArrow("UR", "CTL", "nil, error").
		Deactivate("UR").
		Arrow("CTL", "LR", "FindByLogin(login)").
		Activate("LR").
		Arrow("LR", "DB", "SELECT * FROM librarians WHERE login = $1").
		Activate("DB").
		AltStart("Librarian found and correct password").
		DashedArrow("DB", "LR", "librarian row").
		Deactivate("DB").
		DashedArrow("LR", "CTL", "*Librarian, nil").
		Deactivate("LR").
		DashedArrow("CTL", "C", "200 {id, last_name, first_name, role: librarian}").
		Else("Not found or wrong password").
		DashedArrow("CTL", "C", "401 Incorrect username or password").
		End().
		Blank().
		End().
		Deactivate("CTL")
}

// Diagram 2 — Loan (check then record)

func loanDiagram() *uml.SequenceDiagram {
	return uml.NewSequenceDiagram("Borrowing a copy").
		Actor("C", "Client (logged-in user)").
		Participant("CTL", "LoanController").
		Participant("LR", "LoanRepository").
		Participant("DB", "PostgreSQL").
		Blank().
		Separator("Step 1 - Check").
		Arrow("C", "CTL", "GET /api/loans/check?barcode=X&user_id=Y").
		Activate("CTL").
		Arrow("CTL", "LR", "Check(userId, barcode)").
		Activate("LR").
		Arrow("LR", "DB", "SELECT copy + business checks").
		Activate("DB").
		AltStart("Every check passes").
		DashedArrow("DB", "LR", "LoanPreview").
		Deactivate("DB").
		DashedArrow("LR", "CTL", "LoanPreview, nil").
		Deactivate("LR").
		DashedArrow("CTL", "C", "200 {title, barcode, deposit, current_balance, new_balance}").
		Else("Copy already on loan").
		DashedArrow("CTL", "C", "422 This copy is already on loan").
		Else("Overdue loans").
		DashedArrow("CTL", "C", "422 You have overdue loans").
		Else("Insufficient deposit balance").
		DashedArrow("CTL", "C", "422 Insufficient deposit balance").
		End().
		Deactivate("CTL").
		Blank().
		Separator("Step 2 - Record").
		Arrow("C", "CTL", "POST /api/loans\\n{barcode, user_id}").
		Activate("CTL").
		Arrow("CTL", "LR", "Borrow(userId, barcode)").
		Activate("LR").
		Arrow("LR", "DB", "UPDATE copies SET is_borrowed=true, borrower_id, dates").
		Activate("DB").
		Arrow("LR", "DB", "UPDATE users SET deposit_balance -= deposit").
		Deactivate("DB").
		DashedArrow("LR", "CTL", "nil").
		Deactivate("LR").
		DashedArrow("CTL", "C", "201 {message: Loan registered successfully}").
		Deactivate("CTL")
}

// Diagram 3 — Creating an item (book or journal)

func createItemDiagram() *uml.SequenceDiagram {
	return uml.NewSequenceDiagram("Creating an item (book or journal)").
		Actor("C", "Client (librarian)").
		Participant("CTL", "CreateItemController").
		Participant("FAC", "ItemFactory").
		Participant("IR", "ItemRepository").
		Participant("AR", "AuthorRepository").
		Participant("LOR", "LocationRepository").
		Participant("DB", "PostgreSQL").
		Blank().
		Separator("Loading the form").
		Arrow("C", "CTL", "GET /api/authors").
		Activate("CTL").
		Arrow("CTL", "AR", "FindAll()").
		Activate("AR").
		Arrow("AR", "DB", "SELECT id, last_name, first_name FROM authors ORDER BY last_name").
		Activate("DB").
		DashedArrow("DB", "AR", "[]Author").
		Deactivate("DB").
		DashedArrow("AR", "CTL", "[]Author, nil").
		Deactivate("AR").
		DashedArrow("CTL", "C", "200 [authors]").
		Blank().
		Arrow("C", "CTL", "GET /api/locations").
		Arrow("CTL", "LOR", "FindAll()").
		Activate("LOR").
		Arrow("LOR", "DB", "SELECT locations JOIN categories").
		Activate("DB").
		DashedArrow("DB", "LOR", "[]LocationSummary").
		Deactivate("DB").
		DashedArrow("LOR", "CTL", "[]LocationSummary, nil").
		Deactivate("LOR").
		DashedArrow("CTL", "C", "200 [locations]").
		Deactivate("CTL").
		Blank().
		Separator("Submit - Book").
		Arrow("C", "CTL", "POST /api/books\\n{title, isbn, deposit, author_id?, author_last_name?, author_first_name?, location_id}").
		Activate("CTL").
		Note("right of", "CTL", "Validation: title, isbn required\\nlocation_id > 0, deposit >= 0").
		AltStart("New author (author_id = 0)").
		Arrow("CTL", "AR", "Create(&Author{last_name, first_name})").
		Activate("AR").
		Arrow("AR", "DB", "INSERT INTO authors(last_name, first_name) RETURNING id").
		Activate("DB").
		DashedArrow("DB", "AR", "authorId").
		Deactivate("DB").
		DashedArrow("AR", "CTL", "authorId, nil").
		Deactivate("AR").
		End().
		Arrow("CTL", "FAC", "BookFactory{}.CreateItem()").
		Activate("FAC").
		DashedArrow("FAC", "CTL", "*Book (id=0)").
		Deactivate("FAC").
		Arrow("CTL", "IR", "CreateBook(title, deposit, isbn, authorId, locationId)").
		Activate("IR").
		Arrow("IR", "DB", "INSERT INTO books(title, deposit, isbn, author_id, location_id) RETURNING id").
		Activate("DB").
		Note("right of", "DB", "PostgreSQL inheritance:\\nthe row is also visible in items").
		Deactivate("DB").
		DashedArrow("IR", "CTL", "newId, nil").
		Deactivate("IR").
		DashedArrow("CTL", "C", "201 {id, message: Book created successfully}").
		Deactivate("CTL").
		Blank().
		Separator("Submit - Journal").
		Arrow("C", "CTL", "POST /api/journals\\n{title, issue_number, publication_date, deposit, location_id}").
		Activate("CTL").
		Note("right of", "CTL", "Validation: title, issue_number, publication_date required\\nlocation_id > 0").
		Arrow("CTL", "FAC", "JournalFactory{}.CreateItem()").
		Activate("FAC").
		DashedArrow("FAC", "CTL", "*Journal (id=0)").
		Deactivate("FAC").
		Arrow("CTL", "IR", "CreateJournal(title, deposit, issueNumber, publicationDate, locationId)").
		Activate("IR").
		Arrow("IR", "DB", "INSERT INTO journals(title, deposit, issue_number, publication_date, location_id) RETURNING id").
		Activate("DB").
		Deactivate("DB").
		DashedArrow("IR", "CTL", "newId, nil").
		Deactivate("IR").
		DashedArrow("CTL", "C", "201 {id, message: Journal created successfully}").
		Deactivate("CTL")
}

// Diagram 5 — Creating a user

func createUserDiagram() *uml.SequenceDiagram {
	return uml.NewSequenceDiagram("Creating a user").
		Actor("C", "Client (librarian)").
		Participant("CTL", "CreateUserController").
		Participant("DR", "DepartmentRepository").
		Participant("UR", "UserRepository").
		Participant("SR", "StudentRepo | TeacherRepo").
		Participant("DB", "PostgreSQL").
		Blank().
		Separator("Loading the form").
		Arrow("C", "CTL", "GET /api/departments").
		Activate("CTL").
		Arrow("CTL", "DR", "FindAll()").
		Activate("DR").
		Arrow("DR", "DB", "SELECT * FROM departments ORDER BY name").
		Activate("DB").
		DashedArrow("DB", "DR", "[]Department").
		Deactivate("DB").
		DashedArrow("DR", "CTL", "[]Department, nil").
		Deactivate("DR").
		DashedArrow("CTL", "C", "200 [departments]").
		Deactivate("CTL").
		Blank().
		Separator("Submitting the form").
		Arrow("C", "CTL", "POST /api/users\\n{last_name, first_name, phone_number, birth_date,\\nemail, status, study_year?, department_id?, address}").
		Activate("CTL").
		Note("right of", "CTL", "Validation: 10-digit phone number,\\nvalid date, status in {student,teacher,individual}").
		Blank().
		Arrow("CTL", "UR", "LoginExists(login)").
		Activate("UR").
		Arrow("UR", "DB", "SELECT COUNT(*) FROM users/librarians WHERE login = $1").
		Activate("DB").
		DashedArrow("DB", "UR", "count").
		Deactivate("DB").
		DashedArrow("UR", "CTL", "bool, nil").
		Deactivate("UR").
		Blank().
		AltStart("Login already taken").
		DashedArrow("CTL", "C", "409 An account already exists with this login").
		Else("Login available").
		Arrow("CTL", "DB", "INSERT INTO addresses(...) RETURNING id").
		Activate("DB").
		DashedArrow("DB", "CTL", "addressId").
		Deactivate("DB").
		Arrow("CTL", "SR", "Create(*Student | *Teacher | *Individual)").
		Activate("SR").
		Arrow("SR", "DB", "INSERT INTO students|teachers|users(...) RETURNING id").
		Activate("DB").
		DashedArrow("DB", "SR", "newId").
		Deactivate("DB").
		DashedArrow("SR", "CTL", "newId, nil").
		Deactivate("SR").
		DashedArrow("CTL", "C", "201 {login, password, message}").
		End().
		Deactivate("CTL")
}

// Diagram 6 — A user's active loans

func listLoansDiagram() *uml.SequenceDiagram {
	return uml.NewSequenceDiagram("Active loans").
		Actor("C", "Client (user)").
		Participant("CTL", "LoanController").
		Participant("LR", "LoanRepository").
		Participant("DB", "PostgreSQL").
		Blank().
		Arrow("C", "CTL", "GET /api/loans?user_id=X").
		Activate("CTL").
		Note("right of", "CTL", "Validation: user_id > 0").
		Arrow("CTL", "LR", "GetLoans(userId)").
		Activate("LR").
		Arrow("LR", "DB", "SELECT c.*, i.title\\nFROM copies c JOIN items i ON i.id = c.item_id\\nWHERE c.borrower_id = $1 AND c.is_borrowed = TRUE").
		Activate("DB").
		DashedArrow("DB", "LR", "[]LoanSummary").
		Deactivate("DB").
		DashedArrow("LR", "CTL", "[]LoanSummary, nil").
		Deactivate("LR").
		AltStart("Loans found").
		DashedArrow("CTL", "C", "200 [{id, barcode, title, start_date, due_date, overdue}]").
		Else("No active loan").
		DashedArrow("CTL", "C", "204 No Content").
		End().
		Deactivate("CTL")
}

// Diagram 7 — Overdue loans (librarian view)

func listOverdueDiagram() *uml.SequenceDiagram {
	return uml.NewSequenceDiagram("Overdue loans").
		Actor("C", "Client (librarian)").
		Participant("CTL", "LoanController").
		Participant("LR", "LoanRepository").
		Participant("DB", "PostgreSQL").
		Blank().
		Arrow("C", "CTL", "GET /api/loans/overdue").
		Activate("CTL").
		Arrow("CTL", "LR", "GetOverdueLoans()").
		Activate("LR").
		Arrow("LR", "DB", "SELECT c.*, i.title, u.last_name, u.first_name\\nFROM copies c\\nJOIN items i ON i.id = c.item_id\\nJOIN users u ON u.id = c.borrower_id\\nWHERE c.is_borrowed = TRUE\\nAND c.loan_due_date < NOW()\\nORDER BY loan_due_date ASC").
		Activate("DB").
		DashedArrow("DB", "LR", "[]OverdueLoan").
		Deactivate("DB").
		DashedArrow("LR", "CTL", "[]OverdueLoan, nil").
		Deactivate("LR").
		AltStart("Overdue loans found").
		DashedArrow("CTL", "C", "200 [{id, barcode, title, due_date, last_name, first_name, phone_number}]").
		Else("No overdue loan").
		DashedArrow("CTL", "C", "204 No Content").
		End().
		Deactivate("CTL")
}

// Diagram 8 — Finding a user (librarian view)

func findUserDiagram() *uml.SequenceDiagram {
	return uml.NewSequenceDiagram("Finding a user").
		Actor("C", "Client (librarian)").
		Participant("CTL", "LoanController").
		Participant("UR", "UserRepository").
		Participant("DB", "PostgreSQL").
		Blank().
		Arrow("C", "CTL", "GET /api/users/search\\n?last_name=X&first_name=Y&postal_code=Z&phone_number=W").
		Activate("CTL").
		Note("right of", "CTL", "At least one non-empty criterion is required").
		AltStart("No criterion provided").
		DashedArrow("CTL", "C", "400 At least one search criterion is required").
		Else("At least one criterion provided").
		Arrow("CTL", "UR", "SearchUsers(lastName, firstName, postalCode, phoneNumber)").
		Activate("UR").
		Arrow("UR", "DB", "SELECT u.*\\nFROM students | teachers | ONLY users u\\nLEFT JOIN addresses a ON a.id = u.address_id\\nWHERE <criteria> ORDER BY last_name, first_name").
		Activate("DB").
		DashedArrow("DB", "UR", "[]UserSummary").
		Deactivate("DB").
		DashedArrow("UR", "CTL", "[]UserSummary, nil").
		Deactivate("UR").
		DashedArrow("CTL", "C", "200 [{id, last_name, first_name, phone_number, role}]").
		End().
		Deactivate("CTL")
}

// Diagram 9 — Use cases (overview of the system)

func useCaseDiagram() *uml.UseCaseDiagram {
	return uml.NewUseCaseDiagram("Use case diagram - Library").
		Direction("left to right direction").
		Blank().
		// Actors
		Actor("LIB", "Librarian").
		Actor("USR", "User").
		Actor("STU", "Student").
		Actor("TEA", "Teacher").
		Blank().
		// Actor inheritance
		ActorInherits("STU", "USR").
		ActorInherits("TEA", "USR").
		Blank().
		// System boundary
		RectangleStart("Library management").
		UseCase("UC_LI", "log_in").
		UseCase("UC_CU", "create_user").
		UseCase("UC_FU", "find_user").
		UseCase("UC_LO", "list_overdue_loans").
		UseCase("UC_RI", "register_item").
		UseCase("UC_SI", "search_item").
		UseCase("UC_LUL", "list_user_loans").
		UseCase("UC_LML", "list_my_loans").
		UseCase("UC_BI", "borrow_item").
		RectangleEnd().
		Blank().
		// Librarian associations
		Association("LIB", "UC_CU").
		Association("LIB", "UC_FU").
		Association("LIB", "UC_LO").
		Association("LIB", "UC_RI").
		Association("LIB", "UC_SI").
		Association("LIB", "UC_LUL").
		Blank().
		// User associations (logged in)
		Association("USR", "UC_SI").
		Association("USR", "UC_LML").
		Association("USR", "UC_BI").
		Blank().
		// <<include>> relations
		Include("UC_CU", "UC_LI").
		Include("UC_FU", "UC_LI").
		Include("UC_RI", "UC_LI").
		Include("UC_LUL", "UC_FU").
		Include("UC_LML", "UC_LI").
		Include("UC_BI", "UC_LI")
}
