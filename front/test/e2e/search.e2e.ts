/**
 * E2E tests — Search screen
 *
 * Endpoint: GET /api/items[?title=...&author=...&isbn=...]
 *
 * Inserted data:
 *  - 1 book "Clean Code" by Robert Martin, isbn='978-TEST-001'
 *  - 1 available copy
 */

import { expect } from 'chai';
import { after, before, describe, it } from 'mocha';
import { TestContext } from './test-context.js';

const API_URL = process.env['API_URL'] ?? 'http://localhost:8080';
const ENDPOINT = `${API_URL}/api/items`;
const ctx = new TestContext();

describe('[E2E] Search', () => {
  before(async () => {
    await ctx.before();
    await ctx.exec(`
      INSERT INTO categories (name) VALUES ('Test');
      INSERT INTO locations (bay_number, shelf_number, level, category_id)
        VALUES (1, 1, 1, 1);
      INSERT INTO authors (last_name, first_name) VALUES ('Martin', 'Robert');
      INSERT INTO books (title, deposit, isbn, author_id)
        VALUES ('Clean Code', 5.0, '978-TEST-001', 1);
      INSERT INTO copies (is_borrowed, barcode, loan_period_days, item_id)
        VALUES (FALSE, 'EX-TEST', 15, (SELECT id FROM books WHERE isbn = '978-TEST-001'));
    `);
  });

  after(async () => await ctx.after());

  // ─── FindAll (no parameter) ──────────────────────────────────────────────────

  describe('[E2E] Search — initial load (no parameter)', () => {
    let status: number;
    let body: unknown;

    before(async () => {
      const res = await fetch(ENDPOINT);
      status = res.status;
      body = status === 200 ? ((await res.json()) as unknown[]) : null;
    });

    it('returns HTTP 200 or 204', () => expect(status).to.be.oneOf([200, 204]));

    it('if 200, returns a non-empty array', () => {
      if (status === 200) {
        expect(body).to.be.an('array').with.length.greaterThan(0);
      }
    });

    it("if 200, every item has a 'title'", () => {
      if (status === 200) {
        for (const item of body as Record<string, unknown>[]) {
          expect(item['title']).to.be.a('string').and.not.empty;
        }
      }
    });
  });

  // ─── Search by title ─────────────────────────────────────────────────────────

  describe('[E2E] Search — by title "Clean"', () => {
    let status: number;
    let body: Record<string, unknown>[];

    before(async () => {
      const res = await fetch(`${ENDPOINT}?title=Clean`);
      status = res.status;
      body = status === 200 ? ((await res.json()) as Record<string, unknown>[]) : [];
    });

    it('returns HTTP 200 or 204', () => expect(status).to.be.oneOf([200, 204]));

    it('if 200, every result contains "Clean" in its title', () => {
      if (status === 200) {
        for (const item of body) {
          expect((item['title'] as string).toLowerCase()).to.include('clean');
        }
      }
    });
  });

  // ─── available alone → 400 ───────────────────────────────────────────────────

  describe('[E2E] Search — available alone (no other criterion)', () => {
    let status: number;
    let body: Record<string, unknown>;

    before(async () => {
      const res = await fetch(`${ENDPOINT}?available=true`);
      status = res.status;
      body = (await res.json()) as Record<string, unknown>;
    });

    it('returns HTTP 400', () => expect(status).to.equal(400));
    it("returns an 'error' field", () => expect(body['error']).to.be.a('string').and.not.empty);
  });

  // ─── Search by author ────────────────────────────────────────────────────────

  describe('[E2E] Search — by author "Martin"', () => {
    let status: number;

    before(async () => {
      const res = await fetch(`${ENDPOINT}?author=Martin`);
      status = res.status;
    });

    it('returns HTTP 200 or 204', () => expect(status).to.be.oneOf([200, 204]));
  });

  // ─── Search by unknown ISBN ──────────────────────────────────────────────────

  describe('[E2E] Search — by unknown ISBN', () => {
    let status: number;

    before(async () => {
      const res = await fetch(`${ENDPOINT}?isbn=000-0-00-000000-0`);
      status = res.status;
    });

    it('returns HTTP 204 (no result)', () => expect(status).to.equal(204));
  });
});
