# Frontend

## Tech stack

- Angular (standalone components)
- TypeScript
- Bootstrap
- RxJS

## Running locally

```bash
make run-front
```

The application is available at `http://localhost:4200`.  
It expects the backend at `http://localhost:8080/api`.

## Languages

The interface is available in English (default) and French, with a switch in the navigation bar. The choice is remembered in `localStorage` and sent to the API in the `Accept-Language` header. Translations live in `front/src/app/i18n/translations.ts`.

## Tests

### Unit (Mocha)

Tests the pure utility functions and the translation dictionaries, without a backend.

```bash
cd front
npm run test:mocha
```

### E2E (Mocha + real API)

Tests the full scenarios against the running API. Requires the backend and the database.

```bash
cd front
npm run test:e2e
```

Covered scenarios:

- User and librarian login
- Home statistics
- Item search (title, author, ISBN)
- User creation (student, teacher, individual)
- Book and journal creation with copies
- Viewing a user's loans
- Viewing and updating the deposit

## Sections

- [Architecture](./architecture.md)
