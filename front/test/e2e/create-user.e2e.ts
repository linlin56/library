/**
 * E2E tests — Create user screen
 *
 * Endpoint: POST /api/users
 *
 * Inserted data:
 *  - 1 department (id=1) — for the teacher tests
 *  - jmartin (Jean Martin, individual) — for the login conflict test
 *    "Julien Martin" → generateLogin("Martin","Julien") = "jmartin" → 409
 */

import { expect } from 'chai';
import { after, before, describe, it } from 'mocha';
import { TestContext } from './test-context.js';

const API_URL = process.env['API_URL'] ?? 'http://localhost:8080';
const ctx = new TestContext();

describe('[E2E] Create user', () => {
  before(async () => {
    await ctx.before();
    await ctx.exec(`
      INSERT INTO departments (name) VALUES ('Informatique');
      INSERT INTO users (last_name, first_name, login, password, email)
        VALUES ('Martin', 'Jean', 'jmartin', 'password', 'jean.martin@test.fr');
    `);
  });

  after(async () => await ctx.after());

  // ─── GET /api/departments ────────────────────────────────────────────────────

  describe('[E2E] Create user — list of departments', () => {
    let status: number;
    let body: Record<string, unknown>[];

    before(async () => {
      const res = await fetch(`${API_URL}/api/departments`);
      status = res.status;
      body = (await res.json()) as Record<string, unknown>[];
    });

    it('returns HTTP 200', () => expect(status).to.equal(200));
    it('returns a non-empty array', () =>
      expect(body).to.be.an('array').with.length.greaterThan(0));
    it("every department has a 'name'", () => {
      for (const d of body) {
        expect(d['name']).to.be.a('string').and.not.empty;
      }
    });
  });

  // ─── POST — valid student ────────────────────────────────────────────────────

  describe('[E2E] Create user — create a valid student', () => {
    let status: number;
    let body: Record<string, unknown>;

    before(async () => {
      const res = await fetch(`${API_URL}/api/users`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          last_name: 'Tester',
          first_name: 'Student',
          phone_number: '0612345678',
          birth_date: '2000-06-15',
          email: 's.tester@etud.fr',
          status: 'student',
          study_year: 'L3',
          street_number: '12',
          street_name: 'rue de la Paix',
          postal_code: '75001',
          city: 'Paris',
        }),
      });
      status = res.status;
      body = (await res.json()) as Record<string, unknown>;
    });

    it('returns HTTP 201', () => expect(status).to.equal(201));
    it("returns a generated 'login'", () => expect(body['login']).to.be.a('string').and.not.empty);
    it("returns a 'password'", () => expect(body['password']).to.be.a('string').and.not.empty);
    it('returns a confirmation message', () =>
      expect(body['message']).to.be.a('string').and.not.empty);
  });

  // ─── POST — login conflict ───────────────────────────────────────────────────

  describe('[E2E] Create user — login already exists (jmartin)', () => {
    let status: number;
    let body: Record<string, unknown>;

    before(async () => {
      // generateLogin("Martin", "Julien") = 'j' + 'martin' = "jmartin" (already in the DB)
      const res = await fetch(`${API_URL}/api/users`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          last_name: 'Martin',
          first_name: 'Julien',
          phone_number: '0611111111',
          birth_date: '2002-03-15',
          email: 'j.martin2@test.fr',
          status: 'student',
          study_year: 'L3',
          street_number: '1',
          street_name: 'rue test',
          postal_code: '75000',
          city: 'Paris',
        }),
      });
      status = res.status;
      body = (await res.json()) as Record<string, unknown>;
    });

    it('returns HTTP 409 (conflict)', () => expect(status).to.equal(409));
    it("returns an 'existing_login' field", () =>
      expect(body['existing_login']).to.be.a('string').and.not.empty);
  });

  // ─── POST — invalid phone number ─────────────────────────────────────────────

  describe('[E2E] Create user — invalid phone number', () => {
    let status: number;
    let body: Record<string, unknown>;

    before(async () => {
      const res = await fetch(`${API_URL}/api/users`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          last_name: 'Dupont',
          first_name: 'Jean',
          phone_number: '123',
          birth_date: '1990-01-01',
          email: 'j.dupont@test.fr',
          status: 'individual',
          street_number: '1',
          street_name: 'rue test',
          postal_code: '75000',
          city: 'Paris',
        }),
      });
      status = res.status;
      body = (await res.json()) as Record<string, unknown>;
    });

    it('returns HTTP 400', () => expect(status).to.equal(400));
    it('returns an error message', () => expect(body['error']).to.be.a('string').and.not.empty);
  });

  // ─── POST — student without year of study ────────────────────────────────────

  describe('[E2E] Create user — student without year of study', () => {
    let status: number;

    before(async () => {
      const res = await fetch(`${API_URL}/api/users`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          last_name: 'Noyear',
          first_name: 'Student',
          phone_number: '0699999999',
          birth_date: '2001-01-01',
          email: 's.noyear@etud.fr',
          status: 'student',
          street_number: '1',
          street_name: 'rue test',
          postal_code: '75000',
          city: 'Paris',
        }),
      });
      status = res.status;
    });

    it('returns HTTP 400', () => expect(status).to.equal(400));
  });

  // ─── POST — incomplete body ──────────────────────────────────────────────────

  describe('[E2E] Create user — incomplete body', () => {
    let status: number;

    before(async () => {
      const res = await fetch(`${API_URL}/api/users`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ last_name: 'Incomplete' }),
      });
      status = res.status;
    });

    it('returns HTTP 400', () => expect(status).to.equal(400));
  });
});
