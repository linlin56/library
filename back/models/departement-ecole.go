package models

type DepartementEcole struct {
	Id  int    `json:"id"`
	Nom string `json:"nom"`
}

func NewDepartementEcole(id int, nom string) *DepartementEcole {
	return &DepartementEcole{
		Id:  id,
		Nom: nom,
	}
}
