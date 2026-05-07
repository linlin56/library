# Bibliotheque — Documentation

Application de gestion de bibliothèque universitaire. Elle couvre la recherche d'ouvrages, la gestion des emprunts, le suivi des cautions et l'administration des utilisateurs.

## Profils utilisateur

Deux rôles distincts :

- **Utilisateur** (étudiant, enseignant, particulier) : recherche, emprunt, consultation de ses emprunts, solde caution
- **Bibliothécaire** : création d'utilisateurs et d'ouvrages, gestion des exemplaires, consultation des emprunts de tous les utilisateurs, suivi des retards, mise à jour des cautions

## Structure du projet

```
back/       API REST Go
front/      Application Angular
docs/       Documentation (ce dossier)
db/         Scripts SQL (schéma + données initiales)
.github/    Pipelines CI
```

## Démarrage rapide

```bash
make sql-start   # Démarre PostgreSQL
make run-back    # Lance l'API sur http://localhost:8080
make run-front   # Lance l'interface sur http://localhost:4200
```

## Dépôts

Le dépôt principal est hébergé sur GitHub. Il est automatiquement synchronisé vers GitLab à chaque push sur `main` via la GitHub Action `mirror-gitlab.yml`.

- GitHub : https://github.com/linlin56/bibliotheque-gin
- GitLab : https://gitlab.com/lin56/bibliotheque

Le CI GitLab (`.gitlab-ci.yml`) s'exécute sur GitLab après chaque synchronisation et reproduit le même pipeline que GitHub Actions.

## Pipeline CI

Le CI est découpé en six workflows enchaînés via `workflow_call` :

| Workflow | Déclencheur | Contenu |
|----------|-------------|---------|
| `ci.yml` | push / PR sur `main` | Orchestrateur — appelle les autres dans l'ordre |
| `build.yml` | appelé par ci | Build Go + build Angular |
| `tests.yml` | appelé après build | Tests unitaires Go + tests Mocha |
| `e2e.yml` | appelé après tests | Tests e2e Mocha contre l'API réelle |
| `swagger.yml` | appelé après e2e | Génération + vérification de la spec Swagger |
| `docs.yml` | appelé après swagger (main uniquement) | Déploiement GitHub Pages via Zensical |
| `lint.yml` | appelé après docs | Lint Go avec golangci-lint |

## Sections

- [Backend](./back/index.md)
- [Frontend](./front/index.md)
