/**
 * E2E tests — My loans & Loans (librarian) screens
 *
 * Endpoints:
 *  GET  /api/loans?user_id=X
 *  GET  /api/loans/overdue
 *  GET  /api/loans/check?user_id=X&barcode=Y
 *  POST /api/loans
 *
 * Inserted data:
 *  - mdupont (student, no loan initially)
 *  - 1 book + 1 copy EX-TEST
 */

import { expect } from 'chai';
import { after, before, describe, it } from 'mocha';
import { TestContext } from './test-context.js';

const API_URL = process.env['API_URL'] ?? 'http://localhost:8080';
const ctx = new TestContext();

describe('[E2E] Loans', () => {
  before(async () => {
    await ctx.before();
    await ctx.exec(`
      INSERT INTO categories (name) VALUES ('Test');
      INSERT INTO locations (bay_number, shelf_number, level, category_id)
        VALUES (1, 1, 1, 1);
      INSERT INTO authors (last_name, first_name) VALUES ('Author', 'Test');
      INSERT INTO books (title, deposit, isbn, author_id)
        VALUES ('Test Book', 5.0, '978-TEST-001', 1);
      INSERT INTO copies (is_borrowed, barcode, loan_period_days, item_id)
        VALUES (FALSE, 'EX-TEST', 15, (SELECT id FROM books WHERE isbn = '978-TEST-001'));
      INSERT INTO students (last_name, first_name, login, password, email, study_year, phone_number, birth_date, deposit_balance)
        VALUES ('Dupont', 'Marie', 'mdupont', 'password', 'marie.dupont@test.fr', 'L3', '0600000001', '2000-01-01', 20.0);
    `);
  });

  after(async () => await ctx.after());

  // ─── invalid user_id ─────────────────────────────────────────────────────────

  describe('[E2E] Loans — invalid user_id', () => {
    let status: number;
    let body: Record<string, unknown>;

    before(async () => {
      const res = await fetch(`${API_URL}/api/loans?user_id=abc`);
      status = res.status;
      body = (await res.json()) as Record<string, unknown>;
    });

    it('returns HTTP 400', () => expect(status).to.equal(400));
    it("returns an 'error' field", () => expect(body['error']).to.be.a('string').and.not.empty);
  });

  // ─── mdupont without loans ───────────────────────────────────────────────────

  describe('[E2E] Loans — mdupont without loans', () => {
    let userId: number;
    let status: number;

    before(async () => {
      const loginRes = await fetch(`${API_URL}/api/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ login: 'mdupont', password: 'password' }),
      });
      const user = (await loginRes.json()) as { id: number };
      userId = user.id;
      const res = await fetch(`${API_URL}/api/loans?user_id=${userId}`);
      status = res.status;
    });

    it('returns HTTP 200 or 204', () => expect(status).to.be.oneOf([200, 204]));
  });

  // ─── check without parameters ────────────────────────────────────────────────

  describe('[E2E] Loans — check with missing parameters', () => {
    let status: number;

    before(async () => {
      const res = await fetch(`${API_URL}/api/loans/check`);
      status = res.status;
    });

    it('returns HTTP 400', () => expect(status).to.equal(400));
  });

  // ─── full cycle: check + borrow ──────────────────────────────────────────────

  describe('[E2E] Loans — full cycle: check then borrow EX-TEST', () => {
    let userId: number;
    let checkStatus: number;
    let borrowStatus: number;
    let borrowBody: Record<string, unknown>;

    before(async () => {
      const loginRes = await fetch(`${API_URL}/api/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ login: 'mdupont', password: 'password' }),
      });
      const user = (await loginRes.json()) as { id: number };
      userId = user.id;

      const checkRes = await fetch(`${API_URL}/api/loans/check?user_id=${userId}&barcode=EX-TEST`);
      checkStatus = checkRes.status;

      const borrowRes = await fetch(`${API_URL}/api/loans`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ user_id: userId, barcode: 'EX-TEST' }),
      });
      borrowStatus = borrowRes.status;
      borrowBody = (await borrowRes.json()) as Record<string, unknown>;
    });

    it('check returns HTTP 200', () => expect(checkStatus).to.equal(200));
    it('borrow returns HTTP 201', () => expect(borrowStatus).to.equal(201));
    it("borrow returns a confirmation 'message'", () =>
      expect(borrowBody['message']).to.be.a('string').and.not.empty);

    it('borrowing the same copy twice returns HTTP 422', async () => {
      const res = await fetch(`${API_URL}/api/loans`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ user_id: userId, barcode: 'EX-TEST' }),
      });
      expect(res.status).to.equal(422);
    });
  });

  // ─── overdue loans ────────────────────────────────────────────────────────────

  describe('[E2E] Overdue loans — list (librarian view)', () => {
    let status: number;

    before(async () => {
      const res = await fetch(`${API_URL}/api/loans/overdue`);
      status = res.status;
    });

    it('returns HTTP 200 or 204', () => expect(status).to.be.oneOf([200, 204]));
  });
});
