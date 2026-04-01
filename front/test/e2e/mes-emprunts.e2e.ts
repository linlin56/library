/**
 * Tests e2e — Écrans Mes Emprunts & Emprunts (bibliothécaire)
 *
 * Endpoints :
 *  GET  /api/emprunts?utilisateur_id=X
 *  GET  /api/emprunts/retard
 *  GET  /api/emprunts/verifier?utilisateur_id=X&code_barre=Y
 *  POST /api/emprunts
 *
 * Données insérées :
 *  - mdupont (étudiant, sans emprunt initialement)
 *  - 1 livre + 1 exemplaire EX-TEST
 */

import { expect } from 'chai';
import { after, before, describe, it } from 'mocha';
import { TestContext } from './test-context.js';

const API_URL = process.env['API_URL'] ?? 'http://localhost:8080';
const ctx = new TestContext();

describe('[E2E] Emprunts', () => {
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
      INSERT INTO etudiants (nom, prenom, login, mot_de_passe, email, annee_etude, numero_telephone, date_de_naissance, solde_caution)
        VALUES ('Dupont', 'Marie', 'mdupont', 'mdp', 'marie.dupont@test.fr', 'L3', '0600000001', '2000-01-01', 20.0);
    `);
  });

  after(async () => await ctx.after());

  // ─── utilisateur_id invalide ─────────────────────────────────────────────────

  describe('[E2E] Emprunts — utilisateur_id invalide', () => {
    let status: number;
    let body: Record<string, unknown>;

    before(async () => {
      const res = await fetch(`${API_URL}/api/emprunts?utilisateur_id=abc`);
      status = res.status;
      body = (await res.json()) as Record<string, unknown>;
    });

    it('retourne HTTP 400', () => expect(status).to.equal(400));
    it("retourne un champ 'erreur'", () => expect(body['erreur']).to.be.a('string').and.not.empty);
  });

  // ─── mdupont sans emprunts ───────────────────────────────────────────────────

  describe('[E2E] Emprunts — mdupont sans emprunts', () => {
    let utilisateurId: number;
    let status: number;

    before(async () => {
      const loginRes = await fetch(`${API_URL}/api/connexion`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ login: 'mdupont', mot_de_passe: 'mdp' }),
      });
      const user = (await loginRes.json()) as { id: number };
      utilisateurId = user.id;
      const res = await fetch(`${API_URL}/api/emprunts?utilisateur_id=${utilisateurId}`);
      status = res.status;
    });

    it('retourne HTTP 200 ou 204', () => expect(status).to.be.oneOf([200, 204]));
  });

  // ─── vérifier sans paramètres ────────────────────────────────────────────────

  describe('[E2E] Emprunts — vérifier avec paramètres manquants', () => {
    let status: number;

    before(async () => {
      const res = await fetch(`${API_URL}/api/emprunts/verifier`);
      status = res.status;
    });

    it('retourne HTTP 400', () => expect(status).to.equal(400));
  });

  // ─── cycle complet : vérifier + emprunter ────────────────────────────────────

  describe('[E2E] Emprunts — cycle complet : vérifier puis emprunter EX-TEST', () => {
    let utilisateurId: number;
    let verifierStatus: number;
    let emprunterStatus: number;
    let emprunterBody: Record<string, unknown>;

    before(async () => {
      const loginRes = await fetch(`${API_URL}/api/connexion`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ login: 'mdupont', mot_de_passe: 'mdp' }),
      });
      const user = (await loginRes.json()) as { id: number };
      utilisateurId = user.id;

      const verifierRes = await fetch(
        `${API_URL}/api/emprunts/verifier?utilisateur_id=${utilisateurId}&code_barre=EX-TEST`,
      );
      verifierStatus = verifierRes.status;

      const emprunterRes = await fetch(`${API_URL}/api/emprunts`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ utilisateur_id: utilisateurId, code_barre: 'EX-TEST' }),
      });
      emprunterStatus = emprunterRes.status;
      emprunterBody = (await emprunterRes.json()) as Record<string, unknown>;
    });

    it('vérification retourne HTTP 200', () => expect(verifierStatus).to.equal(200));
    it('emprunt retourne HTTP 201', () => expect(emprunterStatus).to.equal(201));
    it("emprunt retourne un 'message' de confirmation", () =>
      expect(emprunterBody['message']).to.be.a('string').and.not.empty);

    it('double emprunt du même exemplaire retourne HTTP 422', async () => {
      const res = await fetch(`${API_URL}/api/emprunts`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ utilisateur_id: utilisateurId, code_barre: 'EX-TEST' }),
      });
      expect(res.status).to.equal(422);
    });
  });

  // ─── emprunts en retard ───────────────────────────────────────────────────────

  describe('[E2E] Emprunts en retard — liste (vue bibliothécaire)', () => {
    let status: number;

    before(async () => {
      const res = await fetch(`${API_URL}/api/emprunts/retard`);
      status = res.status;
    });

    it('retourne HTTP 200 ou 204', () => expect(status).to.be.oneOf([200, 204]));
  });
});
