package models

import "time"

type Exemplaire struct {
	Id                int `gorm:"primaryKey;autoIncrement"`
	EstEmprunte       bool
	DateDebutEmprunt  time.Time
	DateFinEmprunt    time.Time
	DelaiEmpruntJours int
	CodeBarre         string
	EmprunteurId      int
	Emprunteur        *Utilisateur `gorm:"foreignKey:EmprunteurId"`
}

func NewExemplaire(codeBarre string, delaiEmpruntJours int) *Exemplaire {
	return &Exemplaire{
		EstEmprunte:       false,
		DelaiEmpruntJours: delaiEmpruntJours,
		CodeBarre:         codeBarre,
	}
}
