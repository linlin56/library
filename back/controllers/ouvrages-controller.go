package controllers

import (
	"bibliotheque/db"
	"bibliotheque/repositories"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetOuvrages(c *gin.Context) {
titre := c.Query("titre")
auteur := c.Query("auteur")
isbn := c.Query("isbn")
codeBarre := c.Query("code_barre")
codeRevue := c.Query("code_revue")
disponibleStr := c.Query("disponible")
disponible, _ := strconv.ParseBool(disponibleStr)

repo := repositories.NewOuvrageRepository(db.GlobalDBO)

auMoinsUnChamp := titre != "" || auteur != "" || isbn != "" || codeBarre != "" || codeRevue != ""

// Aucun parametre : retourne tous les ouvrages (chargement initial)
if !auMoinsUnChamp && !disponible {
ouvrages, err := repo.FindAll()
if err != nil {
c.JSON(http.StatusInternalServerError, gin.H{"erreur": err.Error()})
return
}
if len(ouvrages) == 0 {
c.Status(http.StatusNoContent)
return
}
c.JSON(http.StatusOK, ouvrages)
return
}

// disponible seul ne suffit pas
if !auMoinsUnChamp {
c.JSON(http.StatusBadRequest, gin.H{"erreur": "Veuillez renseigner au moins un champ."})
return
}

filtres := repositories.FiltresRecherche{
Titre:      titre,
Auteur:     auteur,
Isbn:       isbn,
CodeBarre:  codeBarre,
CodeRevue:  codeRevue,
Disponible: disponible,
}

resultats, err := repo.Rechercher(filtres)
if err != nil {
c.JSON(http.StatusInternalServerError, gin.H{"erreur": err.Error()})
return
}
if len(resultats) == 0 {
c.Status(http.StatusNoContent)
return
}
c.JSON(http.StatusOK, resultats)
}
