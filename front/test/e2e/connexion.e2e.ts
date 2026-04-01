/**
 * Tests e2e — Écran Connexion
 *
 * Données insérées :
 *   - mdupont (Marie Dupont, étudiant, mdp)
 *   - vpetit  (Valérie Petit, bibliothécaire, mdp)
 */

import { expect } from 'chai';
import { before, after, describe, it } from 'mocha';
import { TestContext } from './test-context.js';

const API_URL = process.env['API_URL'] ?? 'http://localhost:8080';
const ENDPOINT = `${API_URL}/api/connexion`;
const ctx = new TestContext();

describe('[E2E] Connexion', () => {
  before(async () => {
    await ctx.before();
    await ctx.exec(`
      INSERT INTO etudiants (nom, prenom, login, mot_de_passe, email, annee_etude, numero_telephone, date_de_naissance, solde_caution)
        VALUES ('Dupont', 'Marie', 'mdupont', 'mdp', 'marie.dupont@test.fr', 'L3', '0600000001', '2000-01-01', 20.0);
      INSERT INTO bibliothecaires (nom, prenom, login, mot_de_passe)
        VALUES ('Petit', 'Valérie', 'vpetit', 'mdp');
    `);
  });

  after(async () => await ctx.after());

  // ─── Cas nominal — étudiant ─────────────────────────────────────────────────

  describe('[E2E] Connexion — étudiant valide (mdupont / mdp)', () => {
    let status: number;
    let body: Record<string, unknown>;

    before(async () => {
      const res = await fetch(ENDPOINT, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ login: 'mdupont', mot_de_passe: 'mdp' }),
      });
      status = res.status;
      body = (await res.json()) as Record<string, unknown>;
    });

    it('retourne HTTP 200', () => expect(status).to.equal(200));
    it('retourne le nom "Dupont"', () => expect(body['nom']).to.equal('Dupont'));
    it('retourne le prénom "Marie"', () => expect(body['prenom']).to.equal('Marie'));
    it('retourne le rôle "utilisateur"', () => expect(body['role']).to.equal('utilisateur'));
    it('retourne un id numérique positif', () =>
      expect(body['id']).to.be.a('number').and.greaterThan(0));
    it('retourne un message de confirmation', () =>
      expect(body['message']).to.be.a('string').and.not.empty);
  });

  // ─── Cas nominal — bibliothécaire ──────────────────────────────────────────

  describe('[E2E] Connexion — bibliothécaire valide (vpetit / mdp)', () => {
    let status: number;
    let body: Record<string, unknown>;

    before(async () => {
      const res = await fetch(ENDPOINT, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ login: 'vpetit', mot_de_passe: 'mdp' }),
      });
      status = res.status;
      body = (await res.json()) as Record<string, unknown>;
    });

    it('retourne HTTP 200', () => expect(status).to.equal(200));
    it('retourne le rôle "bibliothecaire"', () => expect(body['role']).to.equal('bibliothecaire'));
    it('retourne le nom "Petit"', () => expect(body['nom']).to.equal('Petit'));
    it('retourne le prénom "Valérie"', () => expect(body['prenom']).to.equal('Valérie'));
  });

  // ─── Mauvais mot de passe ───────────────────────────────────────────────────

  describe('[E2E] Connexion — mauvais mot de passe', () => {
    let status: number;
    let body: Record<string, unknown>;

    before(async () => {
      const res = await fetch(ENDPOINT, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ login: 'mdupont', mot_de_passe: 'mauvais' }),
      });
      status = res.status;
      body = (await res.json()) as Record<string, unknown>;
    });

    it('retourne HTTP 401', () => expect(status).to.equal(401));
    it("retourne un champ 'erreur'", () => expect(body['erreur']).to.be.a('string').and.not.empty);
  });

  // ─── Login inexistant ───────────────────────────────────────────────────────

  describe('[E2E] Connexion — login inexistant', () => {
    let status: number;

    before(async () => {
      const res = await fetch(ENDPOINT, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ login: 'utilisateur_inconnu', mot_de_passe: 'mdp' }),
      });
      status = res.status;
    });

    it('retourne HTTP 401', () => expect(status).to.equal(401));
  });

  // ─── Body incomplet ─────────────────────────────────────────────────────────

  describe('[E2E] Connexion — body sans mot de passe', () => {
    let status: number;

    before(async () => {
      const res = await fetch(ENDPOINT, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ login: 'mdupont' }),
      });
      status = res.status;
    });

    it('retourne HTTP 400', () => expect(status).to.equal(400));
  });

  describe('[E2E] Connexion — body vide', () => {
    let status: number;

    before(async () => {
      const res = await fetch(ENDPOINT, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({}),
      });
      status = res.status;
    });

    it('retourne HTTP 400', () => expect(status).to.equal(400));
  });
});
