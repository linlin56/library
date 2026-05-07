# Backend

## Stack technique

- Go
- Gin (routeur HTTP)
- PostgreSQL (héritage natif de tables)
- Swagger (swaggo/swag)

## Démarrage local

```bash
make run-back
```

L'API écoute sur `http://localhost:8080/api`.  
L'interface Swagger est disponible sur `http://localhost:8080/swagger/index.html`.

## Variables d'environnement

| Variable | Valeur par défaut | Description |
|----------|------------------|-------------|
| `DB_HOST` | `localhost` | Hôte PostgreSQL |
| `DB_PORT` | `5432` | Port PostgreSQL |
| `DB_USER` | `postgres` | Utilisateur PostgreSQL |
| `DB_PASSWORD` | _(vide)_ | Mot de passe PostgreSQL |
| `DB_NAME` | `bibliotheque` | Nom de la base |
| `DB_SSLMODE` | `disable` | Mode SSL |

## Tests unitaires

```bash
cd back
go test ./test/...
```

Les tests sont dans `back/test/` et ne nécessitent pas de base de données active.

## Sections

- [Architecture](./architecture.md)
- [Endpoints](./swagger.md)
