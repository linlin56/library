// Commande gen-uml : génère les diagrammes de séquence PlantUML du projet
// ainsi qu'un diagramme de classes via goplantuml.
//
// Usage :
//
//	go run ./cmd/gen-uml            → écrit dans back/diagrams/
//	go run ./cmd/gen-uml -out /tmp  → répertoire personnalisé
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"bibliotheque/uml"
)

// Ce fichier sert à générer des diagrammes de séquence via une lib Go.
// Pour comparer avec mes diagrammes faits à la main sur Paradigm et voir si tout correspond.
func main() {
	out := flag.String("out", "diagrams", "Répertoire de sortie pour les fichiers .puml")
	flag.Parse()

	if err := os.MkdirAll(*out, 0o755); err != nil {
		log.Fatalf("Impossible de créer le répertoire de sortie : %v", err)
	}

	diagrams := map[string]uml.Diagram{
		// Liste des diagrammes à générer
		"01_connexion":           connexionDiagram(),
		"02_emprunt":             empruntDiagram(),
		"03_creer_ouvrage":       creerOuvrageDiagram(),
		// 04_recherche.puml est maintenu manuellement — ne pas ecraser.
		// "04_recherche": rechercheDiagram(),
		"05_creer_utilisateur":   creerUtilisateurDiagram(),
		"06_lister_emprunts":     listerEmpruntsDiagram(),
		"07_lister_retards":      listerRetardsDiagram(),
		"08_trouver_utilisateur": trouverUtilisateurDiagram(),
		"09_cas_utilisation":     casDUtilisationDiagram(),
	}

	for name, d := range diagrams {
		// Écrit chaque diagramme dans un fichier .puml
		path := filepath.Join(*out, name+".puml")
		if err := d.WriteFile(path); err != nil {
			log.Fatalf("Erreur écriture %s : %v", path, err)
		}
		fmt.Printf("OK  %s\n", path)
	}

	genClassDiagram(*out)
}

// genClassDiagram invoque goplantuml pour générer le diagramme de classes.
func genClassDiagram(outDir string) {
	path, err := exec.LookPath("goplantuml")
	if err != nil {
		fmt.Println("INFO goplantuml non trouve dans le PATH — diagramme de classes ignore.")
		fmt.Println("     Installez via : go install github.com/jfeliu007/goplantuml@latest")
		return
	}

	outFile := filepath.Join(outDir, "00_classes.puml")
	f, err := os.Create(outFile)
	if err != nil {
		log.Printf("Erreur creation %s : %v", outFile, err)
		return
	}
	defer f.Close()

	cmd := exec.Command(path, "-recursive", ".")
	cmd.Stdout = f
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Printf("goplantuml a retourne une erreur : %v", err)
		return
	}
	fmt.Printf("OK  %s\n", outFile)
}

// Diagramme 1 — Connexion utilisateur

func connexionDiagram() *uml.SequenceDiagram {
	return uml.NewSequenceDiagram("Connexion d'un utilisateur").
		Actor("C", "Client").
		Participant("CTL", "ConnexionController").
		Participant("UR", "UtilisateurRepository").
		Participant("BR", "BibliothecaireRepository").
		Participant("DB", "PostgreSQL").
		Blank().
		Arrow("C", "CTL", "POST /api/connexion\\n{login, mot_de_passe}").
		Activate("CTL").
		Blank().
		Arrow("CTL", "UR", "FindByLogin(login)").
		Activate("UR").
		Arrow("UR", "DB", "SELECT * FROM ONLY utilisateurs WHERE login = $1").
		Activate("DB").
		Blank().
		AltStart("Utilisateur trouve").
		DashedArrow("DB", "UR", "ligne utilisateur").
		Deactivate("DB").
		DashedArrow("UR", "CTL", "*Utilisateur, nil").
		Deactivate("UR").
		AltStart("Mot de passe correct").
		DashedArrow("CTL", "C", "200 {id, nom, prenom, role: utilisateur}").
		Else("Mot de passe incorrect").
		DashedArrow("CTL", "C", "401 Login ou mot de passe incorrect").
		End().
		Blank().
		Else("Utilisateur non trouve").
		DashedArrow("DB", "UR", "sql.ErrNoRows").
		Deactivate("DB").
		DashedArrow("UR", "CTL", "nil, error").
		Deactivate("UR").
		Arrow("CTL", "BR", "FindByLogin(login)").
		Activate("BR").
		Arrow("BR", "DB", "SELECT * FROM bibliothecaires WHERE login = $1").
		Activate("DB").
		AltStart("Bibliothecaire trouve et mdp correct").
		DashedArrow("DB", "BR", "ligne bibliothecaire").
		Deactivate("DB").
		DashedArrow("BR", "CTL", "*Bibliothecaire, nil").
		Deactivate("BR").
		DashedArrow("CTL", "C", "200 {id, nom, prenom, role: bibliothecaire}").
		Else("Non trouve ou mauvais mot de passe").
		DashedArrow("CTL", "C", "401 Login ou mot de passe incorrect").
		End().
		Blank().
		End().
		Deactivate("CTL")
}

// Diagramme 2 — Emprunt (verification puis enregistrement)

func empruntDiagram() *uml.SequenceDiagram {
	return uml.NewSequenceDiagram("Emprunt d'un exemplaire").
		Actor("C", "Client (utilisateur connecte)").
		Participant("CTL", "EmpruntController").
		Participant("ER", "EmpruntRepository").
		Participant("DB", "PostgreSQL").
		Blank().
		Separator("Etape 1 - Verification").
		Arrow("C", "CTL", "GET /api/emprunts/verifier?code_barre=X&utilisateur_id=Y").
		Activate("CTL").
		Arrow("CTL", "ER", "Verifier(codeBarre, utilisateurId)").
		Activate("ER").
		Arrow("ER", "DB", "SELECT exemplaire + verifications metier").
		Activate("DB").
		AltStart("Toutes les verifications passent").
		DashedArrow("DB", "ER", "PreviewEmprunt").
		Deactivate("DB").
		DashedArrow("ER", "CTL", "PreviewEmprunt, nil").
		Deactivate("ER").
		DashedArrow("CTL", "C", "200 {titre, auteur, caution, date_fin_prevue}").
		Else("Exemplaire deja emprunte").
		DashedArrow("CTL", "C", "400 Cet exemplaire est deja emprunte").
		Else("Solde caution insuffisant").
		DashedArrow("CTL", "C", "400 Solde de caution insuffisant").
		Else("Limite de 3 emprunts atteinte").
		DashedArrow("CTL", "C", "400 Limite de 3 emprunts simultanees atteinte").
		End().
		Deactivate("CTL").
		Blank().
		Separator("Etape 2 - Enregistrement").
		Arrow("C", "CTL", "POST /api/emprunts\\n{code_barre, utilisateur_id}").
		Activate("CTL").
		Arrow("CTL", "ER", "Emprunter(codeBarre, utilisateurId)").
		Activate("ER").
		Arrow("ER", "DB", "BEGIN TRANSACTION").
		Activate("DB").
		Arrow("ER", "DB", "UPDATE exemplaires SET est_emprunte=true, emprunteur_id, dates").
		Arrow("ER", "DB", "UPDATE utilisateurs SET solde_caution -= caution").
		Arrow("ER", "DB", "COMMIT").
		Deactivate("DB").
		DashedArrow("ER", "CTL", "nil").
		Deactivate("ER").
		DashedArrow("CTL", "C", "200 {message: Emprunt enregistre avec succes}").
		Deactivate("CTL")
}

// Diagramme 3 — Creation d'un ouvrage (livre ou revue)

func creerOuvrageDiagram() *uml.SequenceDiagram {
	return uml.NewSequenceDiagram("Creation d'un ouvrage (livre ou revue)").
		Actor("C", "Client (bibliothecaire)").
		Participant("CTL", "CreerOuvrageController").
		Participant("FAB", "FabriqueOuvrage").
		Participant("OR", "OuvrageRepository").
		Participant("AR", "AuteurRepository").
		Participant("EMR", "EmplacementRepository").
		Participant("DB", "PostgreSQL").
		Blank().
		Separator("Chargement du formulaire").
		Arrow("C", "CTL", "GET /api/auteurs").
		Activate("CTL").
		Arrow("CTL", "AR", "FindAll()").
		Activate("AR").
		Arrow("AR", "DB", "SELECT id, nom, prenom FROM auteurs ORDER BY nom").
		Activate("DB").
		DashedArrow("DB", "AR", "[]Auteur").
		Deactivate("DB").
		DashedArrow("AR", "CTL", "[]Auteur, nil").
		Deactivate("AR").
		DashedArrow("CTL", "C", "200 [auteurs]").
		Blank().
		Arrow("C", "CTL", "GET /api/emplacements").
		Arrow("CTL", "EMR", "FindAll()").
		Activate("EMR").
		Arrow("EMR", "DB", "SELECT emplacements JOIN categories").
		Activate("DB").
		DashedArrow("DB", "EMR", "[]EmplacementResume").
		Deactivate("DB").
		DashedArrow("EMR", "CTL", "[]EmplacementResume, nil").
		Deactivate("EMR").
		DashedArrow("CTL", "C", "200 [emplacements]").
		Deactivate("CTL").
		Blank().
		Separator("Soumission - Livre").
		Arrow("C", "CTL", "POST /api/livres\\n{titre, isbn, caution, auteur_id?, auteur_nom?, auteur_prenom?, emplacement_id}").
		Activate("CTL").
		Note("right of", "CTL", "Validation : titre, isbn requis\\nemplacement_id > 0, caution >= 0").
		AltStart("Auteur nouveau (auteur_id = 0)").
		Arrow("CTL", "AR", "Create(&Auteur{nom, prenom})").
		Activate("AR").
		Arrow("AR", "DB", "INSERT INTO auteurs(nom, prenom) RETURNING id").
		Activate("DB").
		DashedArrow("DB", "AR", "auteurId").
		Deactivate("DB").
		DashedArrow("AR", "CTL", "auteurId, nil").
		Deactivate("AR").
		End().
		Arrow("CTL", "FAB", "FabriqueLivre{}.CreerOuvrage()").
		Activate("FAB").
		DashedArrow("FAB", "CTL", "*Livre (id=0)").
		Deactivate("FAB").
		Arrow("CTL", "OR", "CreateLivre(titre, caution, isbn, auteurId, emplacementId)").
		Activate("OR").
		Arrow("OR", "DB", "BEGIN TRANSACTION").
		Activate("DB").
		Arrow("OR", "DB", "INSERT INTO ouvrages(titre, caution, emplacement_id) RETURNING id").
		Arrow("OR", "DB", "INSERT INTO livres(id, isbn, auteur_id)").
		Arrow("OR", "DB", "COMMIT").
		Deactivate("DB").
		DashedArrow("OR", "CTL", "newId, nil").
		Deactivate("OR").
		DashedArrow("CTL", "C", "201 {id, message: Livre cree avec succes}").
		Deactivate("CTL").
		Blank().
		Separator("Soumission - Revue").
		Arrow("C", "CTL", "POST /api/revues\\n{titre, numero, date_parution, caution, emplacement_id}").
		Activate("CTL").
		Note("right of", "CTL", "Validation : titre, numero, date_parution requis\\nemplacement_id > 0").
		Arrow("CTL", "FAB", "FabriqueRevue{}.CreerOuvrage()").
		Activate("FAB").
		DashedArrow("FAB", "CTL", "*Revue (id=0)").
		Deactivate("FAB").
		Arrow("CTL", "OR", "CreateRevue(titre, caution, numero, dateParution, emplacementId)").
		Activate("OR").
		Arrow("OR", "DB", "BEGIN TRANSACTION").
		Activate("DB").
		Arrow("OR", "DB", "INSERT INTO ouvrages(titre, caution, emplacement_id) RETURNING id").
		Arrow("OR", "DB", "INSERT INTO revues(id, numero, date_parution)").
		Arrow("OR", "DB", "COMMIT").
		Deactivate("DB").
		DashedArrow("OR", "CTL", "newId, nil").
		Deactivate("OR").
		DashedArrow("CTL", "C", "201 {id, message: Revue creee avec succes}").
		Deactivate("CTL")
}

// Diagramme 4 — Recherche avancee d'ouvrages

func rechercheDiagram() *uml.SequenceDiagram {
	return uml.NewSequenceDiagram("Recherche avancee d'ouvrages").
		Actor("C", "Client").
		Participant("CTL", "OuvragesController").
		Participant("OR", "OuvrageRepository").
		Participant("DB", "PostgreSQL").
		Blank().
		Arrow("C", "CTL", "GET /api/ouvrages?titre=X&auteur=Y&isbn=Z&code_barre=B&disponible=true").
		Activate("CTL").
		Arrow("CTL", "OR", "Rechercher(filtres)").
		Activate("OR").
		Note("right of", "OR", "Construction dynamique\\nde la requete SQL\\nselon les filtres fournis").
		Arrow("OR", "DB",
			"SELECT ouvrages+livres+auteurs UNION ALL ouvrages+revues\\n"+
				"WHERE <filtres> ORDER BY titre").
		Activate("DB").
		DashedArrow("DB", "OR", "[]OuvrageResultat").
		Deactivate("DB").
		DashedArrow("OR", "CTL", "[]OuvrageResultat, nil").
		Deactivate("OR").
		AltStart("Resultats trouves").
		DashedArrow("CTL", "C", "200 [{id, titre, type, auteur, isbn, date_parution, caution, disponibles}]").
		Else("Aucun resultat").
		DashedArrow("CTL", "C", "200 []").
		End().
		Deactivate("CTL")
}

// Diagramme 5 — Creation d'un utilisateur

func creerUtilisateurDiagram() *uml.SequenceDiagram {
	return uml.NewSequenceDiagram("Creation d'un utilisateur").
		Actor("C", "Client (bibliothecaire)").
		Participant("CTL", "CreerUtilisateurController").
		Participant("DR", "DepartementRepository").
		Participant("UR", "UtilisateurRepository").
		Participant("SR", "EtudiantRepo | EnseignantRepo").
		Participant("DB", "PostgreSQL").
		Blank().
		Separator("Chargement du formulaire").
		Arrow("C", "CTL", "GET /api/departements").
		Activate("CTL").
		Arrow("CTL", "DR", "FindAll()").
		Activate("DR").
		Arrow("DR", "DB", "SELECT * FROM departements_ecole ORDER BY nom").
		Activate("DB").
		DashedArrow("DB", "DR", "[]DepartementEcole").
		Deactivate("DB").
		DashedArrow("DR", "CTL", "[]DepartementEcole, nil").
		Deactivate("DR").
		DashedArrow("CTL", "C", "200 [departements]").
		Deactivate("CTL").
		Blank().
		Separator("Soumission du formulaire").
		Arrow("C", "CTL", "POST /api/utilisateurs\\n{nom, prenom, telephone, date_naissance,\\nemail, statut, annee_etude?, departement_id?, adresse}").
		Activate("CTL").
		Note("right of", "CTL", "Validation : telephone 10 chiffres,\\ndate valide, statut in {etudiant,enseignant,particulier}").
		Blank().
		Arrow("CTL", "UR", "LoginExists(login)").
		Activate("UR").
		Arrow("UR", "DB", "SELECT COUNT(*) FROM utilisateurs WHERE login = $1").
		Activate("DB").
		DashedArrow("DB", "UR", "count").
		Deactivate("DB").
		DashedArrow("UR", "CTL", "bool, nil").
		Deactivate("UR").
		Blank().
		AltStart("Login deja pris").
		DashedArrow("CTL", "C", "409 Un compte existe deja avec ce login").
		Else("Login disponible").
		Arrow("CTL", "DB", "INSERT INTO adresses(...) RETURNING id").
		Activate("DB").
		DashedArrow("DB", "CTL", "adresseId").
		Deactivate("DB").
		Arrow("CTL", "SR", "Create(*Etudiant | *Enseignant | *Particulier)").
		Activate("SR").
		Arrow("SR", "DB", "BEGIN\\nINSERT INTO utilisateurs(...) RETURNING id\\nINSERT INTO etudiants|enseignants(id,...)\\nCOMMIT").
		Activate("DB").
		DashedArrow("DB", "SR", "newId").
		Deactivate("DB").
		DashedArrow("SR", "CTL", "newId, nil").
		Deactivate("SR").
		DashedArrow("CTL", "C", "201 {login, mot_de_passe, message}").
		End().
		Deactivate("CTL")
}

// Diagramme 6 — Liste des emprunts actifs d'un utilisateur

func listerEmpruntsDiagram() *uml.SequenceDiagram {
	return uml.NewSequenceDiagram("Liste des emprunts actifs").
		Actor("C", "Client (utilisateur)").
		Participant("CTL", "EmpruntController").
		Participant("ER", "EmpruntRepository").
		Participant("DB", "PostgreSQL").
		Blank().
		Arrow("C", "CTL", "GET /api/emprunts?utilisateur_id=X").
		Activate("CTL").
		Note("right of", "CTL", "Validation : utilisateur_id > 0").
		Arrow("CTL", "ER", "GetEmprunts(utilisateurId)").
		Activate("ER").
		Arrow("ER", "DB", "SELECT e.*, o.titre, o.caution\\nFROM exemplaires e JOIN ouvrages o ON o.id = e.ouvrage_id\\nWHERE e.emprunteur_id = $1 AND e.est_emprunte = TRUE").
		Activate("DB").
		DashedArrow("DB", "ER", "[]EmpruntItem").
		Deactivate("DB").
		DashedArrow("ER", "CTL", "[]EmpruntItem, nil").
		Deactivate("ER").
		AltStart("Emprunts trouves").
		DashedArrow("CTL", "C", "200 [{code_barre, titre, date_debut, date_fin_prevue, caution}]").
		Else("Aucun emprunt actif").
		DashedArrow("CTL", "C", "200 []").
		End().
		Deactivate("CTL")
}

// Diagramme 7 — Emprunts en retard (vue bibliothecaire)

func listerRetardsDiagram() *uml.SequenceDiagram {
	return uml.NewSequenceDiagram("Liste des emprunts en retard").
		Actor("C", "Client (bibliothecaire)").
		Participant("CTL", "EmpruntController").
		Participant("ER", "EmpruntRepository").
		Participant("DB", "PostgreSQL").
		Blank().
		Arrow("C", "CTL", "GET /api/emprunts/retard").
		Activate("CTL").
		Arrow("CTL", "ER", "GetEmpruntsEnRetard()").
		Activate("ER").
		Arrow("ER", "DB", "SELECT e.*, o.titre, u.nom, u.prenom\\nFROM exemplaires e\\nJOIN ouvrages o ON o.id = e.ouvrage_id\\nJOIN utilisateurs u ON u.id = e.emprunteur_id\\nWHERE e.est_emprunte = TRUE\\nAND e.date_fin_emprunt < NOW()\\nORDER BY date_fin_emprunt ASC").
		Activate("DB").
		DashedArrow("DB", "ER", "[]EmpruntRetardItem").
		Deactivate("DB").
		DashedArrow("ER", "CTL", "[]EmpruntRetardItem, nil").
		Deactivate("ER").
		AltStart("Retards trouves").
		DashedArrow("CTL", "C", "200 [{code_barre, titre, emprunteur, date_fin_prevue, nb_jours_retard}]").
		Else("Aucun retard").
		DashedArrow("CTL", "C", "200 []").
		End().
		Deactivate("CTL")
}

// Diagramme 8 — Recherche d'un utilisateur (vue bibliothecaire)

func trouverUtilisateurDiagram() *uml.SequenceDiagram {
	return uml.NewSequenceDiagram("Recherche d'un utilisateur").
		Actor("C", "Client (bibliothecaire)").
		Participant("CTL", "EmpruntController").
		Participant("UR", "UtilisateurRepository").
		Participant("DB", "PostgreSQL").
		Blank().
		Arrow("C", "CTL", "GET /api/utilisateurs/rechercher\\n?nom=X&prenom=Y&code_postal=Z&numero_telephone=W").
		Activate("CTL").
		Note("right of", "CTL", "Au moins un critere non vide est requis").
		AltStart("Aucun critere fourni").
		DashedArrow("CTL", "C", "400 Au moins un critere de recherche est requis").
		Else("Au moins un critere fourni").
		Arrow("CTL", "UR", "RechercherUtilisateurs(nom, prenom, codePostal, telephone)").
		Activate("UR").
		Note("right of", "UR", "Construction dynamique\\nde la requete SQL\\nselon les criteres fournis").
		Arrow("UR", "DB", "SELECT u.*, a.*\\nFROM utilisateurs u JOIN adresses a ON a.id = u.adresse_id\\nWHERE <criteres> ORDER BY nom, prenom").
		Activate("DB").
		DashedArrow("DB", "UR", "[]UtilisateurResume").
		Deactivate("DB").
		DashedArrow("UR", "CTL", "[]UtilisateurResume, nil").
		Deactivate("UR").
		DashedArrow("CTL", "C", "200 [{id, nom, prenom, login, telephone, solde_caution}]").
		End().
		Deactivate("CTL")
}

// Diagramme 9 — Cas d'utilisation (vue globale du système)

func casDUtilisationDiagram() *uml.UseCaseDiagram {
	return uml.NewUseCaseDiagram("Diagramme de cas d'utilisation - Bibliotheque").
		Direction("left to right direction").
		Blank().
		// Acteurs
		Actor("BIB", "Bibliothecaire").
		Actor("UTI", "Utilisateur").
		Actor("ETU", "Etudiant").
		Actor("ENS", "Enseignant").
		Blank().
		// Héritages acteurs
		ActorInherits("ETU", "UTI").
		ActorInherits("ENS", "UTI").
		Blank().
		// Frontière système
		RectangleStart("Gestion bibliotheque").
		UseCase("UC_CO", "connexion").
		UseCase("UC_CU", "creer_utilisateur").
		UseCase("UC_TU", "trouver_un_utilisateur").
		UseCase("UC_LR", "lister_les_retards").
		UseCase("UC_EO", "enregistrer_ouvrage").
		UseCase("UC_RO", "rechercher_ouvrage").
		UseCase("UC_LEU", "lister_emprunts_utilisateur").
		UseCase("UC_LSE", "lister_ses_emprunts").
		UseCase("UC_EU", "emprunter_un_ouvrage").
		RectangleEnd().
		Blank().
		// Associations Bibliothécaire
		Association("BIB", "UC_CU").
		Association("BIB", "UC_TU").
		Association("BIB", "UC_LR").
		Association("BIB", "UC_EO").
		Association("BIB", "UC_RO").
		Association("BIB", "UC_LEU").
		Blank().
		// Associations Utilisateur (connecte)
		Association("UTI", "UC_RO").
		Association("UTI", "UC_LSE").
		Association("UTI", "UC_EU").
		Blank().
		// Relations <<include>>
		Include("UC_CU", "UC_CO").
		Include("UC_TU", "UC_CO").
		Include("UC_EO", "UC_CO").
		Include("UC_LEU", "UC_TU").
		Include("UC_LSE", "UC_CO").
		Include("UC_EU", "UC_CO")
}
