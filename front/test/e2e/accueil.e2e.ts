/**
 * Tests e2e — Écran Accueil
 *
 * Endpoint : GET /api/accueil
 * Données insérées : 1 livre, 1 exemplaire, 1 étudiant.
 */

import { expect } from 'chai';
import { after, before, describe, it } from 'mocha';
import { TestContext } from './test-context.js';

const API_URL = process.env['API_URL'] ?? 'http://localhost:8080';
const ctx = new TestContext();

describe('[E2E] Accueil', () => {
  before(async () => {
    await ctx.before();
    await ctx.exec(`
      INSERT INTO categories (nom) VALUES ('Test');
      INSERT INTO emplacements (numero_travee, numero_etagere, niveau, categorie_id)
        VALUES (1, 1, 1, 1);
      INSERT INTO auteurs (nom, prenom) VALUES ('Auteur', 'Test');
      INSERT INTO livres (titre, caution, isbn, auteur_id)
        VALUES ('Livre Test', 5.0, '978-TEST-001', 1);
      INSERT INTO exemplaires (est_emprunte, code_barre, delai_emprunt_jours, ouvrage_id)
        VALUES (FALSE, 'EX-TEST', 15, (SELECT id FROM livres WHERE isbn = '978-TEST-001'));
      INSERT INTO etudiants (nom, prenom, login, mot_de_passe, email, annee_etude)
        VALUES ('Dupont', 'Marie', 'mdupont', 'mdp', 'marie.dupont@test.fr', 'L3');
    `);
  });

  after(async () => await ctx.after());

  describe('[E2E] Accueil — statistiques catalogue', () => {
    let status: number;
    let body: Record<string, unknown>;

    before(async () => {
      const res = await fetch(`${API_URL}/api/accueil`);
      status = res.status;
      body = (await res.json()) as Record<string, unknown>;
    });

    it('retourne HTTP 200', () => expect(status).to.equal(200));
    it('retourne un champ nb_livres numérique', () => expect(body['nb_livres']).to.be.a('number'));
    it('retourne un champ nb_revues numérique', () => expect(body['nb_revues']).to.be.a('number'));
    it('retourne un champ nb_exemplaires_disponibles numérique', () =>
      expect(body['nb_exemplaires_disponibles']).to.be.a('number'));
    it('retourne un champ nb_utilisateurs numérique', () =>
      expect(body['nb_utilisateurs']).to.be.a('number'));
    it('nb_livres vaut 1 (données insérées)', () => expect(body['nb_livres']).to.equal(1));
    it('nb_utilisateurs vaut 1 (données insérées)', () =>
      expect(body['nb_utilisateurs']).to.equal(1));
  });
});
