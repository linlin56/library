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

// CreerOuvrageResponse décrit la réponse à la création d'un livre ou d'une revue.
type CreerOuvrageResponse struct {
	ID      int    `json:"id"`
	Message string `json:"message"`
}

// ConflictLoginResponse décrit la réponse 409 quand un login existe déjà.
type ConflictLoginResponse struct {
	Message       string `json:"message"`
	LoginExistant string `json:"login_existant"`
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
