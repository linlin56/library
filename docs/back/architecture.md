# Backend architecture

## Folder structure

```
back/
├── main.go              Entry point: routes, DB init, Swagger
├── controllers/         HTTP layer: validation, JSON responses, status codes
├── repositories/        Data access: SQL queries, result structs
├── models/              Domain models: Go structs mapping the tables
├── builder/             Builder pattern for users, students and teachers
├── factory/             Abstract factory for items (books, journals)
├── i18n/                API messages in English and French
├── db/                  Connection and SQL helpers (custom DBO)
├── docs/                Generated Swagger spec (do not edit by hand)
├── uml/, cmd/gen-uml/   PlantUML diagram generator
├── diagrams/            Generated diagrams
└── test/                Go unit tests
```

## Request lifecycle

1. A Gin route in `main.go` receives the request
2. The controller (`controllers/`) validates the parameters and the body
3. The repository (`repositories/`) runs the SQL
4. The controller returns the JSON with the appropriate HTTP code (200 / 201 / 204 / 4xx / 5xx), with messages in the language requested by the client

## Database access

The project uses a custom DBO (`back/db/dbo.go`) with no ORM.

Main methods:

- `QueryRows` — query returning several rows
- `QueryRow` — query returning a single row
- `Exec` — query with no result
- `ExecReturning` — query with `RETURNING id`
- `WithTx` — runs inside a transaction

On startup, `SeedIfEmpty("db/scripts/insert.sql")` creates and seeds the database if it is empty.

## Relational model

The database uses native PostgreSQL inheritance:

```
items           (parent table)
├── books       INHERITS items  — adds isbn, author_id
└── journals    INHERITS items  — adds issue_number, publication_date

users           (parent table)
├── students    INHERITS users  — adds study_year
└── teachers    INHERITS users  — adds department_id
```

Independent tables: `authors`, `categories`, `locations`, `addresses`, `departments`, `librarians`, `copies`.

The repositories explicitly target the parent or child tables as needed (global reads on `items`, inserts into `books`/`journals`).

## Swagger generation

```bash
cd back
~/go/bin/swag init -g main.go -o docs
```

The CI checks that the committed spec is up to date (`git diff --exit-code -- back/docs`).
