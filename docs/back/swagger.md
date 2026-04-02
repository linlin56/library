# Swagger et endpoints

## Accès Swagger

- API base URL: http://localhost:8080/api

## Groupes d'endpoints

### Accueil

- `GET /api/accueil`

### Authentification

- `POST /api/connexion`

### Ouvrages

- `GET /api/ouvrages`
- `GET /api/auteurs`
- `GET /api/emplacements`
- `POST /api/livres`
- `POST /api/revues`

### Exemplaires

- `GET /api/ouvrages/:id/exemplaires`
- `GET /api/ouvrages/:id/exemplaires/tous`
- `POST /api/ouvrages/:id/exemplaires`

### Utilisateurs

- `GET /api/departements`
- `POST /api/utilisateurs`
- `GET /api/utilisateurs/:id/caution`
- `PUT /api/utilisateurs/:id/caution`
- `GET /api/utilisateurs/rechercher`

### Emprunts

- `GET /api/emprunts/verifier`
- `GET /api/emprunts`
- `GET /api/emprunts/retard`
- `POST /api/emprunts`

## Génération de la spec

Depuis `back`:

```bash
~/go/bin/swag init -g main.go -o docs
```

Fichiers produits:

- `back/docs/docs.go`
- `back/docs/swagger.json`
- `back/docs/swagger.yaml`
