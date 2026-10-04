# Backend

## Tech stack

- Go
- Gin (HTTP router)
- PostgreSQL (native table inheritance)
- Swagger (swaggo/swag)

## Running locally

```bash
make run-back
```

The API listens on `http://localhost:8080/api`.  
The Swagger UI is available at `http://localhost:8080/swagger/index.html`.

## Environment variables

| Variable | Default value | Description |
|----------|---------------|-------------|
| `DB_HOST` | `localhost` | PostgreSQL host |
| `DB_PORT` | `5432` | PostgreSQL port |
| `DB_USER` | `postgres` | PostgreSQL user |
| `DB_PASSWORD` | _(empty)_ | PostgreSQL password |
| `DB_NAME` | `library` | Database name |
| `DB_SSLMODE` | `disable` | SSL mode |
| `CORS_ORIGIN` | `http://localhost:4200` | Origin allowed to call the API |

## Languages

API messages (errors and confirmations) are available in English and French. The language is picked from the `Accept-Language` request header, English by default. The catalogue lives in `back/i18n/messages.go`.

## Unit tests

```bash
cd back
go test ./test/...
```

The tests live in `back/test/` and do not need a running database.

## Sections

- [Architecture](./architecture.md)
- [Endpoints](./swagger.md)
