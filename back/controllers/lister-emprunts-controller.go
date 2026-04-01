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
// @Summary      Lister les emprunts d'un utilisateur
// @Description  Retourne les emprunts actifs pour un utilisateur donne.
// @Tags         Emprunts
// @Produce      json
// @Param        utilisateur_id  query     int  true  "ID utilisateur"
// @Success      200             {array}   repositories.EmpruntItem
// @Success      204             {string}  string  "No Content"
// @Failure      400             {object}  ErrorResponse
// @Failure      500             {object}  ErrorResponse
// @Router       /emprunts [get]
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
	if len(items) == 0 {
		c.Status(http.StatusNoContent)
		return
	}
	c.JSON(http.StatusOK, items)
}

// ListerEmpruntsEnRetard : GET /api/emprunts/retard
// Réservé aux bibliothécaires. Retourne tous les emprunts dont la date de retour est dépassée.
// @Summary      Lister les emprunts en retard
// @Description  Retourne tous les emprunts en retard.
// @Tags         Emprunts
// @Produce      json
// @Success      200  {array}   repositories.EmpruntEnRetardItem
// @Success      204  {string}  string  "No Content"
// @Failure      500  {object}  ErrorResponse
// @Router       /emprunts/retard [get]
func ListerEmpruntsEnRetard(c *gin.Context) {
	repo := repositories.NewEmpruntRepository(db.GlobalDBO)
	items, err := repo.GetEmpruntsEnRetard()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erreur": "Impossible de recuperer les emprunts en retard."})
		return
	}
	if len(items) == 0 {
		c.Status(http.StatusNoContent)
		return
	}
	c.JSON(http.StatusOK, items)
}

// RechercherUtilisateurs : GET /api/utilisateurs/rechercher?nom=&prenom=&code_postal=&numero_telephone=
// Réservé aux bibliothécaires. Au moins un champ non vide est requis.
// @Summary      Rechercher des utilisateurs
// @Description  Recherche des utilisateurs par nom, prenom, code postal ou numero de telephone.
// @Tags         Utilisateurs
// @Produce      json
// @Param        nom               query     string  false  "Nom"
// @Param        prenom            query     string  false  "Prenom"
// @Param        code_postal       query     string  false  "Code postal"
// @Param        numero_telephone  query     string  false  "Numero de telephone"
// @Success      200               {array}   repositories.UtilisateurResume
// @Success      204               {string}  string  "No Content"
// @Failure      400               {object}  ErrorResponse
// @Failure      500               {object}  ErrorResponse
// @Router       /utilisateurs/rechercher [get]
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
	if len(utilisateurs) == 0 {
		c.Status(http.StatusNoContent)
		return
	}
	c.JSON(http.StatusOK, utilisateurs)
}
