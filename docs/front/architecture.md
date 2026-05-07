# Architecture Frontend

## Structure des dossiers

```
front/src/app/
├── accueil/                   Page d'accueil
├── connexion/                 Formulaire de connexion
├── recherche/                 Recherche d'ouvrages + modal d'emprunt
├── mes-emprunts/              Emprunts actifs de l'utilisateur connecté
├── creer-utilisateur/         Formulaire de création utilisateur
├── creer-ouvrage/             Formulaire de création livre ou revue
├── emprunts-bibliothecaire/   Gestion des emprunts côté bibliothécaire
├── emprunts-retard/           Liste des emprunts en retard
├── guards/                    Protections de routes
├── services/                  Appels API
└── shared/                    Composants partagés (navbar)
```

## Routes

| Route | Composant | Accès |
|-------|-----------|-------|
| `/` | AccueilComponent | Public |
| `/connexion` | ConnexionComponent | Non connecté (connexionGuard) |
| `/recherche` | RechercheComponent | Public |
| `/mes-emprunts` | MesEmpruntsComponent | Connecté (authGuard) |
| `/creer-utilisateur` | CreerUtilisateurComponent | Bibliothécaire (bibliothecaireGuard) |
| `/creer-ouvrage` | CreerOuvrageComponent | Bibliothécaire (bibliothecaireGuard) |
| `/emprunts-bibliothecaire` | EmpruntsBibliothecaireComponent | Bibliothécaire (bibliothecaireGuard) |
| `/emprunts-retard` | EmpruntsRetardComponent | Bibliothécaire (bibliothecaireGuard) |
| `**` | — | Redirige vers `/` |

## Guards

- `authGuard` : exige un utilisateur connecté, redirige vers `/connexion` sinon
- `bibliothecaireGuard` : exige le rôle `bibliothecaire`, redirige vers `/` sinon
- `connexionGuard` : bloque l'accès à `/connexion` si déjà connecté, redirige vers `/`

## Gestion de session

`ConnexionService` est le point central de la session :

- Persiste l'utilisateur connecté dans `sessionStorage`
- Expose un `BehaviorSubject<UtilisateurConnecte | null>` (`utilisateurConnecte$`)
- Méthodes : `connecter()`, `deconnecter()`, `estConnecte()`, `getUtilisateurConnecte()`

## Services API

| Service | Endpoints appelés |
|---------|------------------|
| `AccueilService` | `GET /accueil` |
| `ConnexionService` | `POST /connexion` |
| `RechercheService` | `GET /ouvrages` |
| `EmpruntService` | `GET /emprunts/verifier`, `POST /emprunts`, `DELETE /emprunts/:id`, `GET /emprunts`, `GET /emprunts/retard`, `GET /ouvrages/:id/exemplaires`, `GET /ouvrages/:id/exemplaires/tous`, `POST /ouvrages/:id/exemplaires` |
| `UtilisateurService` | `GET /utilisateurs/rechercher`, `GET /utilisateurs/:id/caution`, `PUT /utilisateurs/:id/caution` |
| `CreerUtilisateurService` | `GET /departements`, `POST /utilisateurs` |
| `CreerOuvrageService` | `GET /auteurs`, `GET /emplacements`, `POST /livres`, `POST /revues` |

## Modèles TypeScript

### Session

```typescript
UtilisateurConnecte  { id, nom, prenom, role, message }
UtilisateurResume    { id, nom, prenom, numero_telephone, role }
```

### Ouvrages

```typescript
Ouvrage     { id, titre, caution, type, isbn?, auteur?, numero?, exemplaires_disponibles, emplacement? }
Auteur      { id, nom, prenom }
Emplacement { id, numero_travee, numero_etagere, niveau, categorie_nom }
Departement { Id, Nom }
```

### Exemplaires

```typescript
ExemplaireDisponible { id, code_barre }
ExemplaireComplet    { id, code_barre, est_emprunte, date_fin_emprunt? }
```

### Emprunts

```typescript
EmpruntItem         { id, code_barre, titre, date_debut, date_fin, en_retard }
PreviewEmprunt      { titre, code_barre, caution, solde_actuel, nouveau_solde }
EmpruntEnRetardItem { id, code_barre, titre, date_fin, nom, prenom, numero_telephone }
```

### Caution et création

```typescript
CautionInfo              { solde_caution, caution_totale }
CreerUtilisateurResponse { login, mot_de_passe, message }
CreerOuvrageResponse     { id, message }
FiltresRecherche         { titre?, auteur?, isbn?, codeBarre?, codeRevue?, disponible? }
```

## Pages — rôles fonctionnels

**Accueil** : affiche les 4 statistiques globales (livres, revues, exemplaires disponibles, utilisateurs).

**Connexion** : formulaire login/mot de passe. Redirige automatiquement si déjà connecté.

**Recherche** : filtrage des ouvrages par titre, auteur, ISBN, code barre, code revue ou disponibilité. Pour un bibliothécaire, affiche aussi le bouton d'emprunt et la gestion des exemplaires dans une modale.

**Mes emprunts** : liste les emprunts actifs de l'utilisateur connecté avec les dates de retour, indicateur de retard et bouton "rendre".

**Creer utilisateur** : formulaire réactif avec champs conditionnels selon le statut (année d'étude pour étudiant, département pour enseignant). Affiche le login et le mot de passe générés après création.

**Creer ouvrage** : création en deux étapes — étape 1 : saisie des informations du livre ou de la revue (avec autocomplétion de l'auteur) ; étape 2 : ajout d'exemplaires avec code-barre et délai d'emprunt.

**Emprunts bibliothécaire** : recherche d'un utilisateur, consultation de ses emprunts actifs et modification du montant de la caution totale.

**Emprunts retard** : liste des emprunts dépassant leur date de retour, avec les coordonnées de l'utilisateur concerné.
