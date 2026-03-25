package models

import "time"

type Exemplaire struct {
	Id                int          `json:"id"`
	EstEmprunte       bool         `json:"est_emprunte"`
	DateDebutEmprunt  time.Time    `json:"date_debut_emprunt"`
	DateFinEmprunt    time.Time    `json:"date_fin_emprunt"`
	DelaiEmpruntJours int          `json:"delai_emprunt_jours"`
	CodeBarre         string       `json:"code_barre"`
	EmprunteurId      int          `json:"emprunteur_id"`
	Emprunteur        *Utilisateur `json:"emprunteur,omitempty"`
}

func NewExemplaire(codeBarre string, delaiEmpruntJours int) *Exemplaire {
	return &Exemplaire{
		EstEmprunte:       false,
		DelaiEmpruntJours: delaiEmpruntJours,
		CodeBarre:         codeBarre,
	}
}
