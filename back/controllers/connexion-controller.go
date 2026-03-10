package controllers

import (
	"bibliotheque/db"
	"bibliotheque/repositories"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ConnexionRequest struct {
	Login      string `json:"login" binding:"required"`
	MotDePasse string `json:"mot_de_passe" binding:"required"`
}

type ConnexionResponse struct {
	Nom    string `json:"nom"`
	Prenom string `json:"prenom"`
	Role   string `json:"role"`
}

func Connexion(c *gin.Context) {
	var req ConnexionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erreur": "Champs requis manquants."})
		return
	}

	utilisateurRepo := repositories.NewUtilisateurRepository(db.GlobalDBO)
	utilisateur, err := utilisateurRepo.FindByLogin(req.Login)
	if err == nil {
		if utilisateur.MotDePasse != req.MotDePasse {
			c.JSON(http.StatusUnauthorized, gin.H{"erreur": "Login ou mot de passe incorrect."})
			return
		}
		c.JSON(http.StatusOK, ConnexionResponse{Nom: utilisateur.Nom, Prenom: utilisateur.Prenom, Role: "utilisateur"})
		return
	}

	biblioRepo := repositories.NewBibliothécaireRepository(db.GlobalDBO)
	bibliothecaire, err2 := biblioRepo.FindByLogin(req.Login)
	if err2 != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"erreur": "Login ou mot de passe incorrect."})
		return
	}
	if bibliothecaire.MotDePasse != req.MotDePasse {
		c.JSON(http.StatusUnauthorized, gin.H{"erreur": "Login ou mot de passe incorrect."})
		return
	}
	c.JSON(http.StatusOK, ConnexionResponse{Nom: bibliothecaire.Nom, Prenom: bibliothecaire.Prenom, Role: "bibliothecaire"})
}
