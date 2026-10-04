/**
 * E2E tests — Home screen
 *
 * Endpoint: GET /api/home
 * Inserted data: 1 book, 1 copy, 1 student.
 */

import { expect } from 'chai';
import { after, before, describe, it } from 'mocha';
import { TestContext } from './test-context.js';

const API_URL = process.env['API_URL'] ?? 'http://localhost:8080';
const ctx = new TestContext();

describe('[E2E] Home', () => {
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
      INSERT INTO students (last_name, first_name, login, password, email, study_year)
        VALUES ('Dupont', 'Marie', 'mdupont', 'password', 'marie.dupont@test.fr', 'L3');
    `);
  });

  after(async () => await ctx.after());

  describe('[E2E] Home — catalogue statistics', () => {
    let status: number;
    let body: Record<string, unknown>;

    before(async () => {
      const res = await fetch(`${API_URL}/api/home`);
      status = res.status;
      body = (await res.json()) as Record<string, unknown>;
    });

    it('returns HTTP 200', () => expect(status).to.equal(200));
    it('returns a numeric book_count field', () => expect(body['book_count']).to.be.a('number'));
    it('returns a numeric journal_count field', () =>
      expect(body['journal_count']).to.be.a('number'));
    it('returns a numeric available_copy_count field', () =>
      expect(body['available_copy_count']).to.be.a('number'));
    it('returns a numeric user_count field', () => expect(body['user_count']).to.be.a('number'));
    it('book_count is 1 (inserted data)', () => expect(body['book_count']).to.equal(1));
    it('user_count is 1 (inserted data)', () => expect(body['user_count']).to.equal(1));
  });
});
