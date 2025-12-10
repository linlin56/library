package models

import "time"

type Exemplaire struct {
	EstEmprunte       bool
	DateDebutEmprunt  time.Time
	DateFinEmprunt    time.Time
	DelaiEmpruntJours int
	CodeBarre         string
	Emprunteur        *Utilisateur
}

func NewExemplaire(codeBarre string, delaiEmpruntJours int) *Exemplaire {
	return &Exemplaire{
		EstEmprunte:       false,
		DelaiEmpruntJours: delaiEmpruntJours,
		CodeBarre:         codeBarre,
	}
}
