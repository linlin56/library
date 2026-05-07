# Frontend

## Stack technique

- Angular (standalone components)
- TypeScript
- Bootstrap
- RxJS

## Démarrage local

```bash
make run-front
```

L'application est disponible sur `http://localhost:4200`.  
Elle attend le backend sur `http://localhost:8080/api`.

## Tests

### Unitaires (Mocha)

Teste le rendu du composant racine sans backend.

```bash
cd front
npm run test:mocha
```

### E2E (Mocha + API réelle)

Teste les scénarios complets contre l'API en cours d'exécution. Nécessite le backend et la base de données.

```bash
cd front
npm run test:e2e
```

Scénarios couverts :

- Connexion utilisateur et bibliothécaire
- Statistiques d'accueil
- Recherche d'ouvrages (titre, auteur, ISBN)
- Création utilisateur (étudiant, enseignant, particulier)
- Création livre et revue avec exemplaires
- Consultation des emprunts d'un utilisateur
- Consultation et mise à jour de la caution

## Sections

- [Architecture](./architecture.md)
