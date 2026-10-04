# Library — Documentation

University library management application. It covers item search, loan management, deposit tracking and user administration. The interface is available in English (default) and French.

## User profiles

Two distinct roles:

- **User** (student, teacher, individual): search, borrow, view their loans and deposit balance
- **Librarian**: create users and items, manage copies, view every user's loans, track overdue loans, update deposits

## Project structure

```
back/       Go REST API
front/      Angular application
docs/       Documentation (this folder)
db/         SQL scripts (schema + initial data)
.github/    CI pipelines
```

## Quick start

**Development mode** (Go + Angular with live reload):

```bash
make sql-start   # Starts PostgreSQL
make run-back    # Starts the API on http://localhost:8080
make run-front   # Starts the UI on http://localhost:4200
```

**Docker mode — local build**:

```bash
make docker-start              # Builds the images locally and starts the 3 containers
make docker-stop               # Stops and removes the containers
```

**Docker mode — published images**:

```bash
make docker-start-from-github  # Pulls from GitHub Container Registry (ghcr.io) and starts
make docker-start-from-gitlab  # Pulls from GitLab Container Registry and starts
```

Both commands trigger an interactive `docker login` before the pull. A Personal Access Token with the `read:packages` scope is required.

- Front end: http://localhost
- API: http://localhost:8080
- The database is persisted in a Docker volume (`postgres_data`).

## Repositories

The main repository is hosted on GitHub. It is automatically mirrored to GitLab on every push to `main` by the `mirror-gitlab.yml` GitHub Action.

- GitHub: https://github.com/linlin56/library
- GitLab: https://gitlab.com/lin56/library

The GitLab CI (`.gitlab-ci.yml`) runs on GitLab after each sync and reproduces the same pipeline as GitHub Actions.

## CI pipeline

The CI is split into workflows chained with `workflow_call`:

| Workflow | Trigger | Content |
|----------|---------|---------|
| `ci.yml` | push / PR on `main` | Orchestrator — calls the others in order |
| `build.yml` | called by ci | Go build + Angular build |
| `tests.yml` | called after build | Go unit tests + Mocha tests |
| `e2e.yml` | called after tests | Mocha e2e tests against the real API |
| `swagger.yml` | called after e2e | Generates + checks the Swagger spec |
| `docs.yml` | called after swagger (main only) | GitHub Pages deployment with Zensical |
| `lint.yml` | called after docs | Go lint with golangci-lint |
| `docker.yml` | after a successful CI on `main` | Builds and pushes the images to GHCR |

## Sections

- [Backend](./back/index.md)
- [Frontend](./front/index.md)
