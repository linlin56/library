package controllers

// ErrorResponse décrit une erreur API standard.
type ErrorResponse struct {
	Erreur string `json:"erreur"`
}

// MessageResponse décrit une réponse simple avec message.
type MessageResponse struct {
	Message string `json:"message"`
}

// IDResponse décrit une réponse qui retourne un identifiant.
type IDResponse struct {
	ID int `json:"id"`
}

// UpdateCautionRequest décrit la payload de mise à jour de caution.
type UpdateCautionRequest struct {
	CautionTotale float64 `json:"caution_totale"`
}

// CreerExemplaireRequest décrit la payload de création d'exemplaire.
type CreerExemplaireRequest struct {
	CodeBarre         string `json:"code_barre"`
	DelaiEmpruntJours int    `json:"delai_emprunt_jours"`
}
