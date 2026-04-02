# Architecture Frontend

## Vue d'ensemble

Le frontend est une application Angular structurée par fonctionnalités.

Fichiers d'entrée principaux:

- `front/src/main.ts`
- `front/src/app/app.config.ts`
- `front/src/app/app.routes.ts`

## Routing

Routes principales:

- `/` accueil
- `/connexion`
- `/recherche`
- `/mes-emprunts`
- `/creer-utilisateur`
- `/creer-ouvrage`
- `/emprunts-bibliothecaire`
- `/emprunts-retard`

### Guards utilisés

- `authGuard`: exige un utilisateur connecté
- `bibliothecaireGuard`: exige le rôle `bibliothecaire`
- `connexionGuard`: bloque la page de connexion si déjà connecté

## Gestion de session

Le service `ConnexionService`:

- stocke l'utilisateur connecté dans `sessionStorage`
- expose un `BehaviorSubject` (`utilisateurConnecte$`)
- gère `connecter()`, `deconnecter()`, `estConnecte()`

## Services API

Le frontend s'appuie sur des services Angular dédiés:

- `AccueilService`: statistiques d'accueil
- `ConnexionService`: authentification et session
- `RechercheService`: recherche d'ouvrages
- `EmpruntService`: vérification, emprunts, exemplaires
- `CreerOuvrageService`: auteurs, emplacements, création livre/revue
- `CreerUtilisateurService`: départements, création utilisateur
- `UtilisateurService`: recherche utilisateur, consultation/mise à jour caution

La plupart des endpoints ciblent `http://localhost:8080/api/...`.

## Organisation des composants

Le dossier `front/src/app` est organisé par écrans:

- `accueil`
- `connexion`
- `recherche`
- `mes-emprunts`
- `creer-ouvrage`
- `creer-utilisateur`
- `emprunts-bibliothecaire`
- `emprunts-retard`
- `shared` pour les éléments partagés (ex: navbar)
