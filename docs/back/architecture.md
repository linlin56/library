# Architecture Backend

## Structure des dossiers

```
back/
├── main.go              Point d'entrée : routes, init BDD, Swagger
├── controllers/         Couche HTTP : validation, réponses JSON, codes statut
├── repositories/        Accès données : requêtes SQL, structs de résultat
├── models/              Modèles métier : structs Go mappant les tables
├── db/                  Connexion et utilitaires SQL (DBO maison)
├── docs/                Spec Swagger générée (ne pas modifier manuellement)
└── test/                Tests unitaires Go
```

## Cycle d'une requête

1. Route Gin dans `main.go` reçoit la requête
2. Le contrôleur (`controllers/`) valide les paramètres et le body
3. Le repository (`repositories/`) exécute le SQL
4. Le contrôleur retourne le JSON avec le code HTTP approprié (200 / 201 / 204 / 4xx / 5xx)

## Accès base de données

Le projet utilise un DBO maison (`back/db/dbo.go`) sans ORM.

Méthodes principales :

- `QueryRows` — requête retournant plusieurs lignes
- `QueryRow` — requête retournant une seule ligne
- `Exec` — requête sans retour
- `ExecReturning` — requête avec `RETURNING id`
- `WithTx` — exécution dans une transaction

Au démarrage, `SeedIfEmpty("db/scripts/insert.sql")` crée et peuple la base si elle est vide.

## Modèle relationnel

La base exploite l'héritage natif PostgreSQL :

```
ouvrages        (table parente)
├── livres      INHERITS ouvrages  — ajoute isbn, auteur_id
└── revues      INHERITS ouvrages  — ajoute numero, date_parution

utilisateurs    (table parente)
├── etudiants   INHERITS utilisateurs  — ajoute annee_etude
└── enseignants INHERITS utilisateurs  — ajoute departement_id
```

Tables indépendantes : `auteurs`, `categories`, `emplacements`, `adresses`, `departements_ecole`, `bibliothecaires`, `exemplaires`, `emprunts`.

Les repositories ciblent explicitement les tables parentes ou enfants selon le besoin (lecture globale sur `ouvrages`, création sur `livres`/`revues`).

## Génération Swagger

```bash
cd back
~/go/bin/swag init -g main.go -o docs
```

Le CI vérifie que la spec committée est à jour (`git diff --exit-code -- back/docs`).
