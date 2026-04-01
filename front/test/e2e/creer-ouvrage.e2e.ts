/**
 * Tests e2e — Écran Créer un ouvrage
 *
 * Endpoints :
 *  GET  /api/auteurs
 *  GET  /api/emplacements
 *  POST /api/livres
 *  POST /api/revues
 *
 * Données insérées : 1 catégorie, 1 emplacement (id=1), 1 auteur (id=1).
 */

import { expect } from 'chai';
import { before, after, describe, it } from 'mocha';
import { TestContext } from './test-context.js';

const API_URL = process.env['API_URL'] ?? 'http://localhost:8080';
const ctx = new TestContext();

describe('[E2E] Créer ouvrage', () => {
  before(async () => {
    await ctx.before();
    await ctx.exec(`
      INSERT INTO categories (nom) VALUES ('Test');
      INSERT INTO emplacements (numero_travee, numero_etagere, niveau, categorie_id)
        VALUES (1, 1, 1, 1);
      INSERT INTO auteurs (nom, prenom) VALUES ('Auteur', 'Test');
    `);
  });

  after(async () => await ctx.after());

  // ─── GET /api/auteurs ────────────────────────────────────────────────────────

  describe('[E2E] Créer ouvrage — liste des auteurs', () => {
    let status: number;
    let body: Record<string, unknown>[];

    before(async () => {
      const res = await fetch(`${API_URL}/api/auteurs`);
      status = res.status;
      body = (await res.json()) as Record<string, unknown>[];
    });

    it('retourne HTTP 200', () => expect(status).to.equal(200));
    it('retourne un tableau', () => expect(body).to.be.an('array'));
    it('chaque auteur possède nom et prénom', () => {
      for (const a of body) {
        expect(a['nom']).to.be.a('string').and.not.empty;
        expect(a['prenom']).to.be.a('string');
      }
    });
  });

  // ─── GET /api/emplacements ───────────────────────────────────────────────────

  describe('[E2E] Créer ouvrage — liste des emplacements', () => {
    let status: number;
    let body: Record<string, unknown>[];

    before(async () => {
      const res = await fetch(`${API_URL}/api/emplacements`);
      status = res.status;
      body = (await res.json()) as Record<string, unknown>[];
    });

    it('retourne HTTP 200', () => expect(status).to.equal(200));
    it('retourne un tableau non vide', () =>
      expect(body).to.be.an('array').with.length.greaterThan(0));
  });

  // ─── POST /api/livres — cas nominal ─────────────────────────────────────────

  describe('[E2E] Créer ouvrage — créer un livre valide', () => {
    let status: number;
    let body: Record<string, unknown>;

    before(async () => {
      const res = await fetch(`${API_URL}/api/livres`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          titre: 'Clean Code — test e2e',
          caution: 5.0,
          isbn: '978-0-13-110362-7',
          auteur_nom: 'Martin',
          auteur_prenom: 'Robert',
          emplacement_id: 1,
        }),
      });
      status = res.status;
      body = (await res.json()) as Record<string, unknown>;
    });

    it('retourne HTTP 201', () => expect(status).to.equal(201));
    it("retourne un champ 'id'", () => expect(body['id']).to.be.a('number').and.greaterThan(0));
  });

  // ─── POST /api/livres — sans titre ──────────────────────────────────────────

  describe('[E2E] Créer ouvrage — livre sans titre', () => {
    let status: number;

    before(async () => {
      const res = await fetch(`${API_URL}/api/livres`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          titre: '',
          isbn: '978-0-13-110362-7',
          auteur_nom: 'Martin',
          auteur_prenom: 'Robert',
          emplacement_id: 1,
        }),
      });
      status = res.status;
    });

    it('retourne HTTP 400', () => expect(status).to.equal(400));
  });

  // ─── POST /api/livres — caution négative ────────────────────────────────────

  describe('[E2E] Créer ouvrage — livre avec caution négative', () => {
    let status: number;
    let body: Record<string, unknown>;

    before(async () => {
      const res = await fetch(`${API_URL}/api/livres`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          titre: 'Test caution',
          isbn: '978-0-00-000000-0',
          caution: -10,
          auteur_id: 1,
          emplacement_id: 1,
        }),
      });
      status = res.status;
      body = (await res.json()) as Record<string, unknown>;
    });

    it('retourne HTTP 400', () => expect(status).to.equal(400));
    it("retourne un champ 'erreur'", () => expect(body['erreur']).to.be.a('string').and.not.empty);
  });

  // ─── POST /api/revues — cas nominal ─────────────────────────────────────────

  describe('[E2E] Créer ouvrage — créer une revue valide', () => {
    let status: number;
    let body: Record<string, unknown>;

    before(async () => {
      const res = await fetch(`${API_URL}/api/revues`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          titre: 'Science & Vie — test e2e',
          caution: 2.0,
          numero: 1342,
          date_parution: '2024-01-15',
          emplacement_id: 1,
        }),
      });
      status = res.status;
      body = (await res.json()) as Record<string, unknown>;
    });

    it('retourne HTTP 201', () => expect(status).to.equal(201));
    it("retourne un champ 'id'", () => expect(body['id']).to.be.a('number').and.greaterThan(0));
  });
});
