package controllers

import (
	"bibliotheque/db"
	"bibliotheque/models"
	"bibliotheque/repositories"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var phoneRegexp = regexp.MustCompile(`^0[1-9][0-9]{8}$`)

type CreerUtilisateurRequest struct {
	Nom             string `json:"nom" binding:"required"`
	Prenom          string `json:"prenom" binding:"required"`
	NumeroTelephone string `json:"numero_telephone" binding:"required"`
	DateNaissance   string `json:"date_naissance" binding:"required"`
	Email           string `json:"email" binding:"required"`
	Statut          string `json:"statut" binding:"required,oneof=etudiant enseignant particulier"`
	AnneeEtude      string `json:"annee_etude"`
	DepartementId   int    `json:"departement_id"`
	NumeroRue       string `json:"numero_rue" binding:"required"`
	NomRue          string `json:"nom_rue" binding:"required"`
	CodePostal      string `json:"code_postal" binding:"required"`
	Ville           string `json:"ville" binding:"required"`
}

type CreerUtilisateurResponse struct {
	Login     string `json:"login"`
	MotDePasse string `json:"mot_de_passe"`
	Message   string `json:"message"`
}

func GetDepartements(c *gin.Context) {
	repo := repositories.NewDepartementEcoleRepository(db.GlobalDBO)
	departements, err := repo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erreur": err.Error()})
		return
	}
	c.JSON(http.StatusOK, departements)
}

func CreerUtilisateur(c *gin.Context) {
	var req CreerUtilisateurRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erreur": "Champs requis manquants ou invalides."})
		return
	}

	req.Nom = strings.TrimSpace(req.Nom)
	req.Prenom = strings.TrimSpace(req.Prenom)
	req.Email = strings.TrimSpace(req.Email)
	req.NumeroTelephone = strings.TrimSpace(req.NumeroTelephone)
	req.NomRue = capitaliser(strings.TrimSpace(req.NomRue))
	req.Ville = capitaliser(strings.TrimSpace(req.Ville))

	if !phoneRegexp.MatchString(req.NumeroTelephone) {
		c.JSON(http.StatusBadRequest, gin.H{"erreur": "Le numéro de téléphone doit contenir 10 chiffres et commencer par 0 (ex: 0601020304)."})
		return
	}

	dateNaissance, err := time.Parse("2006-01-02", req.DateNaissance)
	// Normalement ça ne devrait pas arriver car on gère le format côté front
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erreur": "Format de date invalide (attendu : AAAA-MM-JJ)."})
		return
	}

	if req.Statut == "etudiant" && strings.TrimSpace(req.AnneeEtude) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"erreur": "L'année d'étude est requise pour un étudiant."})
		return
	}
	if req.Statut == "enseignant" && req.DepartementId == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"erreur": "Le département est requis pour un enseignant."})
		return
	}

	login := genererLogin(req.Nom, req.Prenom)
	// Pour simplifier le mot de passe sera toujours mdp.
	// En situation réelle, il faudrait bien sûr générer un mdp aléatoire et le hasher
	const motDePasse = "mdp"

	utilisateurRepo := repositories.NewUtilisateurRepository(db.GlobalDBO)
	exists, err := utilisateurRepo.LoginExists(login)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erreur": "Erreur lors de la vérification du compte."})
		return
	}
	if exists {
		c.JSON(http.StatusConflict, gin.H{
			"message":        fmt.Sprintf("Un compte existe déjà avec le login « %s ».", login),
			"login_existant": login,
		})
		return
	}

	var adresseId int
	if err := db.GlobalDBO.ExecReturning(`
		INSERT INTO adresses (ville, code_postal, nom_rue, numero_rue)
		VALUES ($1, $2, $3, $4) RETURNING id`,
		req.Ville, req.CodePostal, req.NomRue, req.NumeroRue,
	).Scan(&adresseId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erreur": "Erreur lors de l'enregistrement de l'adresse."})
		return
	}

	switch req.Statut {
	case "etudiant":
		et := &models.Etudiant{
			Utilisateur: models.Utilisateur{
				Nom:             req.Nom,
				Prenom:          req.Prenom,
				NumeroTelephone: req.NumeroTelephone,
				SoldeCaution:    20.0,
				Email:           req.Email,
				Login:           login,
				MotDePasse:      motDePasse,
				DateDeNaissance: dateNaissance,
				AdresseId:       &adresseId,
			},
			AnneeEtude: req.AnneeEtude,
		}
		etRepo := repositories.NewEtudiantRepository(db.GlobalDBO)
		if _, err := etRepo.Create(et); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"erreur": "Erreur lors de la création de l'étudiant."})
			return
		}

	case "enseignant":
		en := &models.Enseignant{
			Utilisateur: models.Utilisateur{
				Nom:             req.Nom,
				Prenom:          req.Prenom,
				NumeroTelephone: req.NumeroTelephone,
				SoldeCaution:    20.0,
				Email:           req.Email,
				Login:           login,
				MotDePasse:      motDePasse,
				DateDeNaissance: dateNaissance,
				AdresseId:       &adresseId,
			},
			DepartementId: req.DepartementId,
		}
		enRepo := repositories.NewEnseignantRepository(db.GlobalDBO)
		if _, err := enRepo.Create(en); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"erreur": "Erreur lors de la création de l'enseignant."})
			return
		}

	case "particulier":
		var newID int
		err := db.GlobalDBO.ExecReturning(`
			INSERT INTO utilisateurs (nom, prenom, numero_telephone, solde_caution, login, mot_de_passe, date_de_naissance, email, adresse_id)
			VALUES ($1, $2, $3, 20, $4, $5, $6, $7, $8) RETURNING id`,
			req.Nom, req.Prenom, req.NumeroTelephone,
			login, motDePasse, dateNaissance, req.Email, adresseId,
		).Scan(&newID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"erreur": "Erreur lors de la création du particulier."})
			return
		}
	}

	c.JSON(http.StatusCreated, CreerUtilisateurResponse{
		Login:      login,
		MotDePasse: motDePasse,
		Message:    "Utilisateur créé avec succès.",
	})
}

// --- helpers ---
// Pour simplifier on ne stocke pas d'accent
func normaliserChaine(s string) string {
	r := strings.NewReplacer(
		"à", "a", "â", "a", "ä", "a",
		"é", "e", "è", "e", "ê", "e", "ë", "e",
		"î", "i", "ï", "i",
		"ô", "o", "ö", "o",
		"ù", "u", "û", "u", "ü", "u",
		"ç", "c",
		"À", "A", "Â", "A", "Ä", "A",
		"É", "E", "È", "E", "Ê", "E", "Ë", "E",
		"Î", "I", "Ï", "I",
		"Ô", "O", "Ö", "O",
		"Ù", "U", "Û", "U", "Ü", "U",
		"Ç", "C",
	)
	return r.Replace(s)
}

func genererLogin(nom, prenom string) string {
	nomNorm := strings.ToLower(strings.ReplaceAll(normaliserChaine(nom), " ", ""))
	prenomNorm := strings.ToLower(strings.ReplaceAll(normaliserChaine(prenom), " ", ""))
	if len(prenomNorm) == 0 {
		return nomNorm
	}
	return string(prenomNorm[0]) + nomNorm
}

func capitaliser(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(string(s[0])) + strings.ToLower(s[1:])
}
