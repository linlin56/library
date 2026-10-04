package models

// Note (applies to every model)
// The `json` tag defines the snake_case field names used when the API sends JSON to the front end
type Author struct {
	Id        int    `json:"id"`
	LastName  string `json:"last_name"`
	FirstName string `json:"first_name"`
}

func NewAuthor(id int, lastName string, firstName string) *Author {
	return &Author{
		Id:        id,
		LastName:  lastName,
		FirstName: firstName,
	}
}
