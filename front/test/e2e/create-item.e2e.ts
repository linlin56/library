/**
 * E2E tests — Create item screen
 *
 * Endpoints:
 *  GET  /api/authors
 *  GET  /api/locations
 *  POST /api/books
 *  POST /api/journals
 *
 * Inserted data: 1 category, 1 location (id=1), 1 author (id=1).
 */

import { expect } from 'chai';
import { after, before, describe, it } from 'mocha';
import { TestContext } from './test-context.js';

const API_URL = process.env['API_URL'] ?? 'http://localhost:8080';
const ctx = new TestContext();

describe('[E2E] Create item', () => {
  before(async () => {
    await ctx.before();
    await ctx.exec(`
      INSERT INTO categories (name) VALUES ('Test');
      INSERT INTO locations (bay_number, shelf_number, level, category_id)
        VALUES (1, 1, 1, 1);
      INSERT INTO authors (last_name, first_name) VALUES ('Author', 'Test');
    `);
  });

  after(async () => await ctx.after());

  // ─── GET /api/authors ────────────────────────────────────────────────────────

  describe('[E2E] Create item — list of authors', () => {
    let status: number;
    let body: Record<string, unknown>[];

    before(async () => {
      const res = await fetch(`${API_URL}/api/authors`);
      status = res.status;
      body = (await res.json()) as Record<string, unknown>[];
    });

    it('returns HTTP 200', () => expect(status).to.equal(200));
    it('returns an array', () => expect(body).to.be.an('array'));
    it('every author has a last and first name', () => {
      for (const a of body) {
        expect(a['last_name']).to.be.a('string').and.not.empty;
        expect(a['first_name']).to.be.a('string');
      }
    });
  });

  // ─── GET /api/locations ──────────────────────────────────────────────────────

  describe('[E2E] Create item — list of locations', () => {
    let status: number;
    let body: Record<string, unknown>[];

    before(async () => {
      const res = await fetch(`${API_URL}/api/locations`);
      status = res.status;
      body = (await res.json()) as Record<string, unknown>[];
    });

    it('returns HTTP 200', () => expect(status).to.equal(200));
    it('returns a non-empty array', () =>
      expect(body).to.be.an('array').with.length.greaterThan(0));
  });

  // ─── POST /api/books — nominal case ─────────────────────────────────────────

  describe('[E2E] Create item — create a valid book', () => {
    let status: number;
    let body: Record<string, unknown>;

    before(async () => {
      const res = await fetch(`${API_URL}/api/books`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          title: 'Clean Code — e2e test',
          deposit: 5.0,
          isbn: '978-0-13-110362-7',
          author_last_name: 'Martin',
          author_first_name: 'Robert',
          location_id: 1,
        }),
      });
      status = res.status;
      body = (await res.json()) as Record<string, unknown>;
    });

    it('returns HTTP 201', () => expect(status).to.equal(201));
    it("returns an 'id' field", () => expect(body['id']).to.be.a('number').and.greaterThan(0));
  });

  // ─── POST /api/books — without title ────────────────────────────────────────

  describe('[E2E] Create item — book without title', () => {
    let status: number;

    before(async () => {
      const res = await fetch(`${API_URL}/api/books`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          title: '',
          isbn: '978-0-13-110362-7',
          author_last_name: 'Martin',
          author_first_name: 'Robert',
          location_id: 1,
        }),
      });
      status = res.status;
    });

    it('returns HTTP 400', () => expect(status).to.equal(400));
  });

  // ─── POST /api/books — negative deposit ─────────────────────────────────────

  describe('[E2E] Create item — book with a negative deposit', () => {
    let status: number;
    let body: Record<string, unknown>;

    before(async () => {
      const res = await fetch(`${API_URL}/api/books`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          title: 'Deposit test',
          isbn: '978-0-00-000000-0',
          deposit: -10,
          author_id: 1,
          location_id: 1,
        }),
      });
      status = res.status;
      body = (await res.json()) as Record<string, unknown>;
    });

    it('returns HTTP 400', () => expect(status).to.equal(400));
    it("returns an 'error' field", () => expect(body['error']).to.be.a('string').and.not.empty);
  });

  // ─── POST /api/journals — nominal case ──────────────────────────────────────

  describe('[E2E] Create item — create a valid journal', () => {
    let status: number;
    let body: Record<string, unknown>;

    before(async () => {
      const res = await fetch(`${API_URL}/api/journals`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          title: 'Science & Vie — e2e test',
          deposit: 2.0,
          issue_number: 1342,
          publication_date: '2024-01-15',
          location_id: 1,
        }),
      });
      status = res.status;
      body = (await res.json()) as Record<string, unknown>;
    });

    it('returns HTTP 201', () => expect(status).to.equal(201));
    it("returns an 'id' field", () => expect(body['id']).to.be.a('number').and.greaterThan(0));
  });
});
