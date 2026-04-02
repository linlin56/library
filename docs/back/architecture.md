# Architecture Backend

## Vue d'ensemble

Le backend est une API REST en Go, construite avec Gin.
L'entrée applicative est `back/main.go`, qui:

- initialise la connexion BDD (`db.Init()`)
- gère l'API
- expose Swagger

## Structure technique

- `back/controllers`: couche HTTP (validation des paramètres, réponses JSON, codes HTTP)
- `back/repositories`: logique d'accès aux données SQL
- `back/models`: modèles métier
- `back/db`: gestion concrète du SQL

## Accès base de données

Le projet utilise un DBO maison (`back/db/dbo.go`) au lieu d'un ORM.

Points clés:

- `QueryRows`, `QueryRow`, `Exec`, `ExecReturning`
- gestion transactionnelle via `WithTx`
- peuplement automatique via `SeedIfEmpty("db/scripts/insert.sql")`

## Modèle relationnel PostgreSQL

La base s'appuie sur l'héritage PostgreSQL natif:

- `livres INHERITS ouvrages`
- `revues INHERITS ouvrages`
- `etudiants INHERITS utilisateurs`
- `enseignants INHERITS utilisateurs`

Ce choix est cohérent avec les repositories qui ciblent explicitement les tables parentes/enfants selon les besoins.

## Cycle d'une requête

1. Route Gin dans `main.go`
2. Contrôleur (`controllers/*`) valide la requête
3. Repository exécute SQL
4. Contrôleur renvoie JSON (200/201/204/4xx/5xx)
