# Library

University library management application: item search, loans, deposit tracking and user administration.
The interface is available in English (default) and French.

- **Back end**: Go, Gin, PostgreSQL (native table inheritance), Swagger
- **Front end**: Angular, Bootstrap

Documentation: https://linlin56.github.io/library/

## Getting started

Requirements: Go, Node.js 22, and Docker for the database container.

```bash
make sql-start                # PostgreSQL container (port 5435)
make run-back                 # API on http://localhost:8080 (creates and seeds the database on first run)
cd front && npm install       # first run only
make run-front                # UI on http://localhost:4200
```

Or run everything with Docker: `make docker-start` (front end on http://localhost).

## Demo accounts

Every account uses the password `password`.

| Login | Role |
|-------|------|
| `vpetit` | Librarian |
| `jmartin` | User (individual) |
| `mdupont` | User (student) |
| `proche` | User (teacher) |
