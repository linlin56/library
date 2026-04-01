/**
 * Tests e2e — Écran Caution (vue bibliothécaire)
 *
 * Endpoints :
 *  GET /api/utilisateurs/:id/caution
 *  PUT /api/utilisateurs/:id/caution
 *  GET /api/utilisateurs/rechercher?nom=...
 *
 * Données insérées :
 *  - jmartin (Jean Martin, particulier, solde=20, total=20)
 *  - proche  (Pierre Roche, enseignant)  — emprunte EX-TEST dans le test caution < solde
 */

import { expect } from 'chai';
import { before, after, describe, it } from 'mocha';
import { TestContext } from './test-context.js';

const API_URL = process.env['API_URL'] ?? 'http://localhost:8080';
const ctx = new TestContext();

async function getIdByLogin(login: string): Promise<number> {
  const res = await fetch(`${API_URL}/api/connexion`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ login, mot_de_passe: 'mdp' }),
  });
  const body = (await res.json()) as { id: number };
  return body.id;
}

describe('[E2E] Caution', () => {
  before(async () => {
    await ctx.before();
    await ctx.exec(`
      INSERT INTO categories (nom) VALUES ('Test');
      INSERT INTO emplacements (numero_travee, numero_etagere, niveau, categorie_id)
        VALUES (1, 1, 1, 1);
      INSERT INTO departements_ecole (nom) VALUES ('Informatique');
      INSERT INTO auteurs (nom, prenom) VALUES ('Auteur', 'Test');
      INSERT INTO livres (titre, caution, isbn, auteur_id)
        VALUES ('Livre Test', 5.0, '978-TEST-001', 1);
      INSERT INTO exemplaires (est_emprunte, code_barre, delai_emprunt_jours, ouvrage_id)
        VALUES (FALSE, 'EX-TEST', 15, (SELECT id FROM livres WHERE isbn = '978-TEST-001'));
      INSERT INTO utilisateurs (nom, prenom, login, mot_de_passe, email, numero_telephone, date_de_naissance)
        VALUES ('Martin', 'Jean', 'jmartin', 'mdp', 'jean.martin@test.fr', '0600000001', '1990-01-01');
      INSERT INTO enseignants (nom, prenom, login, mot_de_passe, email, departement_id, numero_telephone, date_de_naissance, solde_caution)
        VALUES ('Roche', 'Pierre', 'proche', 'mdp', 'pierre.roche@test.fr', 1, '0600000002', '1985-05-15', 20.0);
    `);
  });

  after(async () => await ctx.after());

  // ─── GET /api/utilisateurs/:id/caution ──────────────────────────────────────

  describe('[E2E] Caution — consulter la caution de jmartin', () => {
    let status: number;
    let body: Record<string, unknown>;

    before(async () => {
      const id = await getIdByLogin('jmartin');
      const res = await fetch(`${API_URL}/api/utilisateurs/${id}/caution`);
      status = res.status;
      body = (await res.json()) as Record<string, unknown>;
    });

    it('retourne HTTP 200', () => expect(status).to.equal(200));
    it("retourne un champ 'solde_caution' numérique", () =>
      expect(body['solde_caution']).to.be.a('number'));
    it("retourne un champ 'caution_totale' numérique", () =>
      expect(body['caution_totale']).to.be.a('number'));
  });

  // ─── GET /api/utilisateurs/:id/caution — id invalide ────────────────────────

  describe('[E2E] Caution — id utilisateur invalide', () => {
    let status: number;

    before(async () => {
      const res = await fetch(`${API_URL}/api/utilisateurs/99999/caution`);
      status = res.status;
    });

    it('retourne HTTP 404', () => expect(status).to.equal(404));
  });

  // ─── PUT /api/utilisateurs/:id/caution — mise à jour valide ─────────────────

  describe('[E2E] Caution — modifier la caution totale de jmartin', () => {
    let status: number;
    let body: Record<string, unknown>;

    before(async () => {
      const id = await getIdByLogin('jmartin');
      const res = await fetch(`${API_URL}/api/utilisateurs/${id}/caution`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ caution_totale: 30.0 }),
      });
      status = res.status;
      body = (await res.json()) as Record<string, unknown>;
    });

    it('retourne HTTP 200', () => expect(status).to.equal(200));
    it("retourne le nouveau 'caution_totale'", () => expect(body['caution_totale']).to.equal(30.0));
  });

  // ─── PUT — montant inférieur au solde emprunté ───────────────────────────────

  describe('[E2E] Caution — caution totale inférieure au montant emprunté (proche)', () => {
    let status: number;
    let body: Record<string, unknown>;

    before(async () => {
      // proche emprunte EX-TEST (caution livre = 5 €)
      // → son solde_caution passe de 20 à 15, montant_emprunté = 5
      // → fixer caution_totale = 1 < 5 doit retourner 400
      const id = await getIdByLogin('proche');
      await fetch(`${API_URL}/api/emprunts`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ utilisateur_id: id, code_barre: 'EX-TEST' }),
      });
      const res = await fetch(`${API_URL}/api/utilisateurs/${id}/caution`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ caution_totale: 1.0 }),
      });
      status = res.status;
      body = (await res.json()) as Record<string, unknown>;
    });

    it('retourne HTTP 400', () => expect(status).to.equal(400));
    it("retourne un champ 'erreur'", () => expect(body['erreur']).to.be.a('string').and.not.empty);
  });

  // ─── GET /api/utilisateurs/rechercher ────────────────────────────────────────

  describe('[E2E] Caution — rechercher utilisateur par nom "Martin"', () => {
    let status: number;
    let body: Record<string, unknown>[];

    before(async () => {
      const res = await fetch(`${API_URL}/api/utilisateurs/rechercher?nom=Martin`);
      status = res.status;
      body = status === 200 ? ((await res.json()) as Record<string, unknown>[]) : [];
    });

    it('retourne HTTP 200 ou 204', () => expect(status).to.be.oneOf([200, 204]));

    it('si 200, retourne au moins un résultat', () => {
      if (status === 200) {
        expect(body).to.be.an('array').with.length.greaterThan(0);
      }
    });

    it("si 200, chaque résultat possède 'nom' et 'prenom'", () => {
      if (status === 200) {
        for (const u of body) {
          expect(u['nom']).to.be.a('string').and.not.empty;
          expect(u['prenom']).to.be.a('string');
        }
      }
    });
  });
});
