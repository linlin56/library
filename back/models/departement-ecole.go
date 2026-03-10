package models

type DepartementEcole struct {
	Id  int `gorm:"primaryKey;autoIncrement"`
	Nom string
}

func NewDepartementEcole(id int, nom string) *DepartementEcole {
	return &DepartementEcole{
		Id:  id,
		Nom: nom,
	}
}
