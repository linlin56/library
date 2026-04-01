/**
 * Tests e2e — Écran Recherche
 *
 * Endpoint : GET /api/ouvrages[?titre=...&auteur=...&isbn=...]
 *
 * Données insérées :
 *  - 1 livre "Clean Code" de Robert Martin, isbn='978-TEST-001'
 *  - 1 exemplaire disponible
 */

import { expect } from 'chai';
import { after, before, describe, it } from 'mocha';
import { TestContext } from './test-context.js';

const API_URL = process.env['API_URL'] ?? 'http://localhost:8080';
const ENDPOINT = `${API_URL}/api/ouvrages`;
const ctx = new TestContext();

describe('[E2E] Recherche', () => {
  before(async () => {
    await ctx.before();
    await ctx.exec(`
      INSERT INTO categories (nom) VALUES ('Test');
      INSERT INTO emplacements (numero_travee, numero_etagere, niveau, categorie_id)
        VALUES (1, 1, 1, 1);
      INSERT INTO auteurs (nom, prenom) VALUES ('Martin', 'Robert');
      INSERT INTO livres (titre, caution, isbn, auteur_id)
        VALUES ('Clean Code', 5.0, '978-TEST-001', 1);
      INSERT INTO exemplaires (est_emprunte, code_barre, delai_emprunt_jours, ouvrage_id)
        VALUES (FALSE, 'EX-TEST', 15, (SELECT id FROM livres WHERE isbn = '978-TEST-001'));
    `);
  });

  after(async () => await ctx.after());

  // ─── FindAll (aucun paramètre) ───────────────────────────────────────────────

  describe('[E2E] Recherche — chargement initial (aucun paramètre)', () => {
    let status: number;
    let body: unknown;

    before(async () => {
      const res = await fetch(ENDPOINT);
      status = res.status;
      body = status === 200 ? ((await res.json()) as unknown[]) : null;
    });

    it('retourne HTTP 200 ou 204', () => expect(status).to.be.oneOf([200, 204]));

    it('si 200, retourne un tableau non vide', () => {
      if (status === 200) {
        expect(body).to.be.an('array').with.length.greaterThan(0);
      }
    });

    it("si 200, chaque ouvrage possède un 'titre'", () => {
      if (status === 200) {
        for (const o of body as Record<string, unknown>[]) {
          expect(o['titre']).to.be.a('string').and.not.empty;
        }
      }
    });
  });

  // ─── Recherche par titre ─────────────────────────────────────────────────────

  describe('[E2E] Recherche — par titre "Clean"', () => {
    let status: number;
    let body: Record<string, unknown>[];

    before(async () => {
      const res = await fetch(`${ENDPOINT}?titre=Clean`);
      status = res.status;
      body = status === 200 ? ((await res.json()) as Record<string, unknown>[]) : [];
    });

    it('retourne HTTP 200 ou 204', () => expect(status).to.be.oneOf([200, 204]));

    it('si 200, tous les résultats contiennent "Clean" dans le titre', () => {
      if (status === 200) {
        for (const o of body) {
          expect((o['titre'] as string).toLowerCase()).to.include('clean');
        }
      }
    });
  });

  // ─── disponible seul → 400 ───────────────────────────────────────────────────

  describe('[E2E] Recherche — disponible seul (sans autre critère)', () => {
    let status: number;
    let body: Record<string, unknown>;

    before(async () => {
      const res = await fetch(`${ENDPOINT}?disponible=true`);
      status = res.status;
      body = (await res.json()) as Record<string, unknown>;
    });

    it('retourne HTTP 400', () => expect(status).to.equal(400));
    it("retourne un champ 'erreur'", () => expect(body['erreur']).to.be.a('string').and.not.empty);
  });

  // ─── Recherche par auteur ────────────────────────────────────────────────────

  describe('[E2E] Recherche — par auteur "Martin"', () => {
    let status: number;

    before(async () => {
      const res = await fetch(`${ENDPOINT}?auteur=Martin`);
      status = res.status;
    });

    it('retourne HTTP 200 ou 204', () => expect(status).to.be.oneOf([200, 204]));
  });

  // ─── Recherche par ISBN inexistant ───────────────────────────────────────────

  describe('[E2E] Recherche — par ISBN inexistant', () => {
    let status: number;

    before(async () => {
      const res = await fetch(`${ENDPOINT}?isbn=000-0-00-000000-0`);
      status = res.status;
    });

    it('retourne HTTP 204 (aucun résultat)', () => expect(status).to.equal(204));
  });
});
