package models

type Department struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func NewDepartment(id int, name string) *Department {
	return &Department{
		Id:   id,
		Name: name,
	}
}
