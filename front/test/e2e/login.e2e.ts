/**
 * E2E tests — Login screen
 *
 * Inserted data:
 *   - mdupont (Marie Dupont, student, password)
 *   - vpetit  (Valérie Petit, librarian, password)
 */

import { expect } from 'chai';
import { after, before, describe, it } from 'mocha';
import { TestContext } from './test-context.js';

const API_URL = process.env['API_URL'] ?? 'http://localhost:8080';
const ENDPOINT = `${API_URL}/api/login`;
const ctx = new TestContext();

describe('[E2E] Login', () => {
  before(async () => {
    await ctx.before();
    await ctx.exec(`
      INSERT INTO students (last_name, first_name, login, password, email, study_year, phone_number, birth_date, deposit_balance)
        VALUES ('Dupont', 'Marie', 'mdupont', 'password', 'marie.dupont@test.fr', 'L3', '0600000001', '2000-01-01', 20.0);
      INSERT INTO librarians (last_name, first_name, login, password)
        VALUES ('Petit', 'Valérie', 'vpetit', 'password');
    `);
  });

  after(async () => await ctx.after());

  // ─── Nominal case — student ─────────────────────────────────────────────────

  describe('[E2E] Login — valid student (mdupont / password)', () => {
    let status: number;
    let body: Record<string, unknown>;

    before(async () => {
      const res = await fetch(ENDPOINT, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ login: 'mdupont', password: 'password' }),
      });
      status = res.status;
      body = (await res.json()) as Record<string, unknown>;
    });

    it('returns HTTP 200', () => expect(status).to.equal(200));
    it('returns the last name "Dupont"', () => expect(body['last_name']).to.equal('Dupont'));
    it('returns the first name "Marie"', () => expect(body['first_name']).to.equal('Marie'));
    it('returns the "user" role', () => expect(body['role']).to.equal('user'));
    it('returns a positive numeric id', () =>
      expect(body['id']).to.be.a('number').and.greaterThan(0));
    it('returns a confirmation message', () =>
      expect(body['message']).to.be.a('string').and.not.empty);
  });

  // ─── Nominal case — librarian ───────────────────────────────────────────────

  describe('[E2E] Login — valid librarian (vpetit / password)', () => {
    let status: number;
    let body: Record<string, unknown>;

    before(async () => {
      const res = await fetch(ENDPOINT, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ login: 'vpetit', password: 'password' }),
      });
      status = res.status;
      body = (await res.json()) as Record<string, unknown>;
    });

    it('returns HTTP 200', () => expect(status).to.equal(200));
    it('returns the "librarian" role', () => expect(body['role']).to.equal('librarian'));
    it('returns the last name "Petit"', () => expect(body['last_name']).to.equal('Petit'));
    it('returns the first name "Valérie"', () => expect(body['first_name']).to.equal('Valérie'));
  });

  // ─── Wrong password ─────────────────────────────────────────────────────────

  describe('[E2E] Login — wrong password', () => {
    let status: number;
    let body: Record<string, unknown>;

    before(async () => {
      const res = await fetch(ENDPOINT, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ login: 'mdupont', password: 'wrong' }),
      });
      status = res.status;
      body = (await res.json()) as Record<string, unknown>;
    });

    it('returns HTTP 401', () => expect(status).to.equal(401));
    it("returns an 'error' field", () => expect(body['error']).to.be.a('string').and.not.empty);
  });

  // ─── Unknown login ──────────────────────────────────────────────────────────

  describe('[E2E] Login — unknown login', () => {
    let status: number;

    before(async () => {
      const res = await fetch(ENDPOINT, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ login: 'unknown_user', password: 'password' }),
      });
      status = res.status;
    });

    it('returns HTTP 401', () => expect(status).to.equal(401));
  });

  // ─── Incomplete body ────────────────────────────────────────────────────────

  describe('[E2E] Login — body without password', () => {
    let status: number;

    before(async () => {
      const res = await fetch(ENDPOINT, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ login: 'mdupont' }),
      });
      status = res.status;
    });

    it('returns HTTP 400', () => expect(status).to.equal(400));
  });

  describe('[E2E] Login — empty body', () => {
    let status: number;

    before(async () => {
      const res = await fetch(ENDPOINT, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({}),
      });
      status = res.status;
    });

    it('returns HTTP 400', () => expect(status).to.equal(400));
  });
});
