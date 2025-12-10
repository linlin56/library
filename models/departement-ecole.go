package models

type DepartementEcole struct {
	Id  int
	Nom string
}

func NewDepartementEcole(id int, nom string) *DepartementEcole {
	return &DepartementEcole{
		Id:  id,
		Nom: nom,
	}
}
