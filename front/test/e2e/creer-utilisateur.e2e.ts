/**
 * Tests e2e — Écran Créer un utilisateur
 *
 * Endpoint : POST /api/utilisateurs
 *
 * Données insérées :
 *  - 1 département (id=1) — pour les tests enseignant
 *  - jmartin (Jean Martin, particulier) — pour le test conflit de login
 *    "Julien Martin" → genererLogin("Martin","Julien") = "jmartin" → 409
 */

import { expect } from 'chai';
import { before, after, describe, it } from 'mocha';
import { TestContext } from './test-context.js';

const API_URL = process.env['API_URL'] ?? 'http://localhost:8080';
const ctx = new TestContext();

describe('[E2E] Créer utilisateur', () => {
  before(async () => {
    await ctx.before();
    await ctx.exec(`
      INSERT INTO departements_ecole (nom) VALUES ('Informatique');
      INSERT INTO utilisateurs (nom, prenom, login, mot_de_passe, email)
        VALUES ('Martin', 'Jean', 'jmartin', 'mdp', 'jean.martin@test.fr');
    `);
  });

  after(async () => await ctx.after());

  // ─── GET /api/departements ───────────────────────────────────────────────────

  describe('[E2E] Créer utilisateur — liste des départements', () => {
    let status: number;
    let body: Record<string, unknown>[];

    before(async () => {
      const res = await fetch(`${API_URL}/api/departements`);
      status = res.status;
      body = (await res.json()) as Record<string, unknown>[];
    });

    it('retourne HTTP 200', () => expect(status).to.equal(200));
    it('retourne un tableau non vide', () =>
      expect(body).to.be.an('array').with.length.greaterThan(0));
    it("chaque département possède un 'nom'", () => {
      for (const d of body) {
        expect(d['nom']).to.be.a('string').and.not.empty;
      }
    });
  });

  // ─── POST — étudiant valide ──────────────────────────────────────────────────

  describe('[E2E] Créer utilisateur — créer un étudiant valide', () => {
    let status: number;
    let body: Record<string, unknown>;

    before(async () => {
      const res = await fetch(`${API_URL}/api/utilisateurs`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          nom: 'Testeur',
          prenom: 'Etudiant',
          numero_telephone: '0612345678',
          date_naissance: '2000-06-15',
          email: 'e.testeur@etud.fr',
          statut: 'etudiant',
          annee_etude: 'L3',
          numero_rue: '12',
          nom_rue: 'rue de la Paix',
          code_postal: '75001',
          ville: 'Paris',
        }),
      });
      status = res.status;
      body = (await res.json()) as Record<string, unknown>;
    });

    it('retourne HTTP 201', () => expect(status).to.equal(201));
    it("retourne un 'login' généré", () => expect(body['login']).to.be.a('string').and.not.empty);
    it("retourne un 'mot_de_passe'", () =>
      expect(body['mot_de_passe']).to.be.a('string').and.not.empty);
    it('retourne un message de confirmation', () =>
      expect(body['message']).to.be.a('string').and.not.empty);
  });

  // ─── POST — conflit login ────────────────────────────────────────────────────

  describe('[E2E] Créer utilisateur — login déjà existant (jmartin)', () => {
    let status: number;
    let body: Record<string, unknown>;

    before(async () => {
      // genererLogin("Martin", "Julien") = 'j' + 'martin' = "jmartin" (déjà en BDD)
      const res = await fetch(`${API_URL}/api/utilisateurs`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          nom: 'Martin',
          prenom: 'Julien',
          numero_telephone: '0611111111',
          date_naissance: '2002-03-15',
          email: 'j.martin2@test.fr',
          statut: 'etudiant',
          annee_etude: 'L3',
          numero_rue: '1',
          nom_rue: 'rue test',
          code_postal: '75000',
          ville: 'Paris',
        }),
      });
      status = res.status;
      body = (await res.json()) as Record<string, unknown>;
    });

    it('retourne HTTP 409 (conflit)', () => expect(status).to.equal(409));
    it("retourne un champ 'login_existant'", () =>
      expect(body['login_existant']).to.be.a('string').and.not.empty);
  });

  // ─── POST — téléphone invalide ───────────────────────────────────────────────

  describe('[E2E] Créer utilisateur — téléphone invalide', () => {
    let status: number;
    let body: Record<string, unknown>;

    before(async () => {
      const res = await fetch(`${API_URL}/api/utilisateurs`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          nom: 'Dupont',
          prenom: 'Jean',
          numero_telephone: '123',
          date_naissance: '1990-01-01',
          email: 'j.dupont@test.fr',
          statut: 'particulier',
          numero_rue: '1',
          nom_rue: 'rue test',
          code_postal: '75000',
          ville: 'Paris',
        }),
      });
      status = res.status;
      body = (await res.json()) as Record<string, unknown>;
    });

    it('retourne HTTP 400', () => expect(status).to.equal(400));
    it("retourne un message d'erreur", () => expect(body['erreur']).to.be.a('string').and.not.empty);
  });

  // ─── POST — étudiant sans année d'étude ──────────────────────────────────────

  describe("[E2E] Créer utilisateur — étudiant sans année d'étude", () => {
    let status: number;

    before(async () => {
      const res = await fetch(`${API_URL}/api/utilisateurs`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          nom: 'Sansan',
          prenom: 'Annee',
          numero_telephone: '0699999999',
          date_naissance: '2001-01-01',
          email: 'a.sansan@etud.fr',
          statut: 'etudiant',
          numero_rue: '1',
          nom_rue: 'rue test',
          code_postal: '75000',
          ville: 'Paris',
        }),
      });
      status = res.status;
    });

    it('retourne HTTP 400', () => expect(status).to.equal(400));
  });

  // ─── POST — body incomplet ───────────────────────────────────────────────────

  describe('[E2E] Créer utilisateur — body incomplet', () => {
    let status: number;

    before(async () => {
      const res = await fetch(`${API_URL}/api/utilisateurs`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ nom: 'Incomplet' }),
      });
      status = res.status;
    });

    it('retourne HTTP 400', () => expect(status).to.equal(400));
  });
});
