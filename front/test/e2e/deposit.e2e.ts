/**
 * E2E tests — Deposit screen (librarian view)
 *
 * Endpoints:
 *  GET /api/users/:id/deposit
 *  PUT /api/users/:id/deposit
 *  GET /api/users/search?last_name=...
 *
 * Inserted data:
 *  - jmartin (Jean Martin, individual, balance=20, total=20)
 *  - proche  (Pierre Roche, teacher) — borrows EX-TEST in the "deposit < balance" test
 */

import { expect } from 'chai';
import { after, before, describe, it } from 'mocha';
import { TestContext } from './test-context.js';

const API_URL = process.env['API_URL'] ?? 'http://localhost:8080';
const ctx = new TestContext();

async function getIdByLogin(login: string): Promise<number> {
  const res = await fetch(`${API_URL}/api/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ login, password: 'password' }),
  });
  const body = (await res.json()) as { id: number };
  return body.id;
}

describe('[E2E] Deposit', () => {
  before(async () => {
    await ctx.before();
    await ctx.exec(`
      INSERT INTO categories (name) VALUES ('Test');
      INSERT INTO locations (bay_number, shelf_number, level, category_id)
        VALUES (1, 1, 1, 1);
      INSERT INTO departments (name) VALUES ('Informatique');
      INSERT INTO authors (last_name, first_name) VALUES ('Author', 'Test');
      INSERT INTO books (title, deposit, isbn, author_id)
        VALUES ('Test Book', 5.0, '978-TEST-001', 1);
      INSERT INTO copies (is_borrowed, barcode, loan_period_days, item_id)
        VALUES (FALSE, 'EX-TEST', 15, (SELECT id FROM books WHERE isbn = '978-TEST-001'));
      INSERT INTO users (last_name, first_name, login, password, email, phone_number, birth_date)
        VALUES ('Martin', 'Jean', 'jmartin', 'password', 'jean.martin@test.fr', '0600000001', '1990-01-01');
      INSERT INTO teachers (last_name, first_name, login, password, email, department_id, phone_number, birth_date, deposit_balance)
        VALUES ('Roche', 'Pierre', 'proche', 'password', 'pierre.roche@test.fr', 1, '0600000002', '1985-05-15', 20.0);
    `);
  });

  after(async () => await ctx.after());

  // ─── GET /api/users/:id/deposit ──────────────────────────────────────────────

  describe("[E2E] Deposit — get jmartin's deposit", () => {
    let status: number;
    let body: Record<string, unknown>;

    before(async () => {
      const id = await getIdByLogin('jmartin');
      const res = await fetch(`${API_URL}/api/users/${id}/deposit`);
      status = res.status;
      body = (await res.json()) as Record<string, unknown>;
    });

    it('returns HTTP 200', () => expect(status).to.equal(200));
    it("returns a numeric 'deposit_balance' field", () =>
      expect(body['deposit_balance']).to.be.a('number'));
    it("returns a numeric 'total_deposit' field", () =>
      expect(body['total_deposit']).to.be.a('number'));
  });

  // ─── GET /api/users/:id/deposit — invalid id ─────────────────────────────────

  describe('[E2E] Deposit — invalid user id', () => {
    let status: number;

    before(async () => {
      const res = await fetch(`${API_URL}/api/users/99999/deposit`);
      status = res.status;
    });

    it('returns HTTP 404', () => expect(status).to.equal(404));
  });

  // ─── PUT /api/users/:id/deposit — valid update ───────────────────────────────

  describe("[E2E] Deposit — update jmartin's total deposit", () => {
    let status: number;
    let body: Record<string, unknown>;

    before(async () => {
      const id = await getIdByLogin('jmartin');
      const res = await fetch(`${API_URL}/api/users/${id}/deposit`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ total_deposit: 30.0 }),
      });
      status = res.status;
      body = (await res.json()) as Record<string, unknown>;
    });

    it('returns HTTP 200', () => expect(status).to.equal(200));
    it("returns the new 'total_deposit'", () => expect(body['total_deposit']).to.equal(30.0));
  });

  // ─── PUT — amount lower than the borrowed amount ─────────────────────────────

  describe('[E2E] Deposit — total deposit lower than the borrowed amount (proche)', () => {
    let status: number;
    let body: Record<string, unknown>;

    before(async () => {
      // proche borrows EX-TEST (book deposit = €5)
      // → their deposit_balance goes from 20 to 15, borrowed amount = 5
      // → setting total_deposit = 1 < 5 must return 400
      const id = await getIdByLogin('proche');
      await fetch(`${API_URL}/api/loans`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ user_id: id, barcode: 'EX-TEST' }),
      });
      const res = await fetch(`${API_URL}/api/users/${id}/deposit`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ total_deposit: 1.0 }),
      });
      status = res.status;
      body = (await res.json()) as Record<string, unknown>;
    });

    it('returns HTTP 400', () => expect(status).to.equal(400));
    it("returns an 'error' field", () => expect(body['error']).to.be.a('string').and.not.empty);
  });

  // ─── GET /api/users/search ───────────────────────────────────────────────────

  describe('[E2E] Deposit — search users by last name "Martin"', () => {
    let status: number;
    let body: Record<string, unknown>[];

    before(async () => {
      const res = await fetch(`${API_URL}/api/users/search?last_name=Martin`);
      status = res.status;
      body = status === 200 ? ((await res.json()) as Record<string, unknown>[]) : [];
    });

    it('returns HTTP 200 or 204', () => expect(status).to.be.oneOf([200, 204]));

    it('if 200, returns at least one result', () => {
      if (status === 200) {
        expect(body).to.be.an('array').with.length.greaterThan(0);
      }
    });

    it("if 200, every result has a 'last_name' and a 'first_name'", () => {
      if (status === 200) {
        for (const u of body) {
          expect(u['last_name']).to.be.a('string').and.not.empty;
          expect(u['first_name']).to.be.a('string');
        }
      }
    });
  });
});
