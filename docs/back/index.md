# Backend

## Stack technique

- Go
- Gin (API HTTP)
- PostgreSQL
- Swagger

## Démarrage local

Depuis la racine du projet:

```bash
make run-back
```

Swagger local:

- http://localhost:8080/api

## Sections

- [Architecture](./architecture.md)
- [Swagger et endpoints](./swagger.md)

## Tests backend

Le backend contient des tests unitaires Go dans `back/test`.

Exécution:

```bash
cd back
go test ./test/...
```
