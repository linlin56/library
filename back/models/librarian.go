package models

type Librarian struct {
	Id        int    `json:"id"`
	LastName  string `json:"last_name"`
	FirstName string `json:"first_name"`
	Login     string `json:"login"`
	Password  string `json:"password"`
}

func NewLibrarian(id int, lastName string, firstName string, login string, password string) *Librarian {
	return &Librarian{
		Id:        id,
		LastName:  lastName,
		FirstName: firstName,
		Login:     login,
		Password:  password,
	}
}
