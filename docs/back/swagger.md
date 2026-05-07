# Endpoints

Base URL : `http://localhost:8080/api`

## Accueil

### GET /accueil

Retourne les statistiques globales de la bibliothèque.

- Accès : public
- Réponse 200 : `{ nb_livres, nb_revues, nb_exemplaires_disponibles, nb_utilisateurs }`

## Authentification

### POST /connexion

Authentifie un utilisateur ou un bibliothécaire.

- Accès : public
- Body : `{ login, mot_de_passe }`
- Réponse 200 : `{ id, nom, prenom, role, message }`
- Réponse 400 : champs manquants
- Réponse 401 : identifiants invalides

## Ouvrages

### GET /ouvrages

Recherche d'ouvrages avec filtres optionnels.

- Accès : public
- Paramètres query (tous optionnels) : `titre`, `auteur`, `isbn`, `code_barre`, `code_revue`, `disponible` (booléen)
- Réponse 200 : tableau d'ouvrages `{ id, titre, caution, type, isbn?, auteur?, numero?, exemplaires_disponibles, emplacement? }`
- Réponse 204 : aucun résultat
- Réponse 400 : paramètre invalide
- Réponse 500 : erreur serveur

### GET /auteurs

Retourne la liste de tous les auteurs.

- Accès : bibliothécaire
- Réponse 200 : tableau `{ id, nom, prenom }`
- Réponse 500 : erreur serveur

### GET /emplacements

Retourne la liste des emplacements disponibles.

- Accès : bibliothécaire
- Réponse 200 : tableau `{ id, numero_travee, numero_etagere, niveau, categorie_nom }`
- Réponse 500 : erreur serveur

### POST /livres

Crée un nouveau livre. L'auteur peut être existant (par `auteur_id`) ou nouveau (par `auteur_nom` + `auteur_prenom`).

- Accès : bibliothécaire
- Body : `{ titre, isbn, caution, emplacement_id, auteur_id?, auteur_nom?, auteur_prenom? }`
- Réponse 201 : `{ id, message }`
- Réponse 400 : données invalides
- Réponse 500 : erreur serveur

### POST /revues

Crée une nouvelle revue.

- Accès : bibliothécaire
- Body : `{ titre, numero, date_parution, caution, emplacement_id }`
- Réponse 201 : `{ id, message }`
- Réponse 400 : données invalides
- Réponse 500 : erreur serveur

## Exemplaires

### GET /ouvrages/:id/exemplaires

Retourne les exemplaires non empruntés d'un ouvrage.

- Accès : public
- Paramètre path : `id` (entier, ID de l'ouvrage)
- Réponse 200 : tableau `{ id, code_barre }`
- Réponse 204 : aucun exemplaire disponible
- Réponse 400 : ID invalide
- Réponse 500 : erreur serveur

### GET /ouvrages/:id/exemplaires/tous

Retourne tous les exemplaires d'un ouvrage (disponibles et empruntés).

- Accès : bibliothécaire
- Paramètre path : `id` (entier)
- Réponse 200 : tableau `{ id, code_barre, est_emprunte, date_fin_emprunt? }`
- Réponse 204 : aucun exemplaire
- Réponse 400 : ID invalide
- Réponse 500 : erreur serveur

### POST /ouvrages/:id/exemplaires

Crée un exemplaire pour un ouvrage donné.

- Accès : bibliothécaire
- Paramètre path : `id` (entier)
- Body : `{ code_barre, delai_emprunt_jours }`
- Réponse 201 : `{ id }`
- Réponse 400 : données invalides

## Utilisateurs

### GET /departements

Retourne la liste des départements d'école.

- Accès : bibliothécaire
- Réponse 200 : tableau `{ id, nom }`
- Réponse 500 : erreur serveur

### POST /utilisateurs

Crée un étudiant, un enseignant ou un particulier.

- Accès : bibliothécaire
- Body : `{ nom, prenom, numero_telephone, date_naissance, email, statut, annee_etude?, departement_id?, numero_rue, nom_rue, code_postal, ville }`
- `statut` : `etudiant` | `enseignant` | `particulier`
- Réponse 201 : `{ login, mot_de_passe, message }`
- Réponse 400 : données invalides ou format incorrect
- Réponse 409 : `{ message, login_existant }` — un compte existe déjà avec ce login
- Réponse 500 : erreur serveur

### GET /utilisateurs/rechercher

Recherche des utilisateurs par critères.

- Accès : bibliothécaire
- Paramètres query (tous optionnels) : `nom`, `prenom`, `code_postal`, `numero_telephone`
- Réponse 200 : tableau `{ id, nom, prenom, numero_telephone, role }`
- Réponse 204 : aucun résultat
- Réponse 400 : paramètre invalide
- Réponse 500 : erreur serveur

### GET /utilisateurs/:id/caution

Retourne le solde de caution d'un utilisateur.

- Accès : bibliothécaire
- Paramètre path : `id` (entier)
- Réponse 200 : `{ solde_caution, caution_totale }`
- Réponse 400 : ID invalide
- Réponse 404 : utilisateur introuvable

### PUT /utilisateurs/:id/caution

Met à jour la caution totale d'un utilisateur et recalcule son solde.

- Accès : bibliothécaire
- Paramètre path : `id` (entier)
- Body : `{ caution_totale }`
- Réponse 200 : `{ solde_caution, caution_totale }`
- Réponse 400 : données invalides

## Emprunts

### GET /emprunts/verifier

Vérifie si un emprunt est possible et retourne un aperçu (impact sur la caution).

- Accès : bibliothécaire
- Paramètres query : `utilisateur_id` (entier), `code_barre` (string)
- Réponse 200 : `{ titre, code_barre, caution, solde_actuel, nouveau_solde }`
- Réponse 400 : paramètre manquant ou invalide
- Réponse 422 : emprunt impossible (caution insuffisante, exemplaire déjà emprunté, etc.)

### POST /emprunts

Enregistre un emprunt.

- Accès : bibliothécaire
- Body : `{ utilisateur_id, code_barre }`
- Réponse 201 : `{ message }`
- Réponse 400 : données invalides
- Réponse 422 : emprunt impossible

### DELETE /emprunts/:id

Retourne un exemplaire et restitue la caution correspondante.

- Accès : bibliothécaire
- Paramètre path : `id` (entier, ID de l'exemplaire)
- Réponse 200 : `{ message }`
- Réponse 400 : ID invalide
- Réponse 422 : exemplaire non emprunté ou introuvable

### GET /emprunts

Retourne les emprunts actifs d'un utilisateur.

- Accès : utilisateur connecté
- Paramètre query : `utilisateur_id` (entier)
- Réponse 200 : tableau `{ id, code_barre, titre, date_debut, date_fin, en_retard }`
- Réponse 204 : aucun emprunt actif
- Réponse 400 : paramètre invalide
- Réponse 500 : erreur serveur

### GET /emprunts/retard

Retourne tous les emprunts dont la date de retour est dépassée.

- Accès : bibliothécaire
- Réponse 200 : tableau `{ id, code_barre, titre, date_fin, nom, prenom, numero_telephone }`
- Réponse 204 : aucun retard
- Réponse 500 : erreur serveur
