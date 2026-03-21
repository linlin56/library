package controllers

import (
	"bibliotheque/db"
	"bibliotheque/repositories"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// ListerEmprunts : GET /api/emprunts?utilisateur_id=X
// Retourne les emprunts actifs d'un utilisateur.
func ListerEmprunts(c *gin.Context) {
	idStr := c.Query("utilisateur_id")
	utilisateurId, err := strconv.Atoi(idStr)
	if err != nil || utilisateurId <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"erreur": "Parametre utilisateur_id invalide."})
		return
	}

	repo := repositories.NewEmpruntRepository(db.GlobalDBO)
	items, err := repo.GetEmprunts(utilisateurId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erreur": "Impossible de recuperer les emprunts."})
		return
	}
	c.JSON(http.StatusOK, items)
}

// RechercherUtilisateurs : GET /api/utilisateurs/rechercher?nom=&prenom=&code_postal=&numero_telephone=
// Réservé aux bibliothécaires. Au moins un champ non vide est requis.
func RechercherUtilisateurs(c *gin.Context) {
	nom := strings.TrimSpace(c.Query("nom"))
	prenom := strings.TrimSpace(c.Query("prenom"))
	codePostal := strings.TrimSpace(c.Query("code_postal"))
	numeroTelephone := strings.TrimSpace(c.Query("numero_telephone"))

	if nom == "" && prenom == "" && codePostal == "" && numeroTelephone == "" {
		c.JSON(http.StatusBadRequest, gin.H{"erreur": "Au moins un critere de recherche est requis."})
		return
	}

	repo := repositories.NewUtilisateurRepository(db.GlobalDBO)
	utilisateurs, err := repo.RechercherUtilisateurs(nom, prenom, codePostal, numeroTelephone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erreur": "Erreur lors de la recherche."})
		return
	}
	c.JSON(http.StatusOK, utilisateurs)
}
