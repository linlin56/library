import { expect } from 'chai';
import { describe, it } from 'mocha';

import { buildSearchParams, isLibrarian, parseUser } from '../src/utils';

// ─── buildSearchParams ────────────────────────────────────────────────────────

describe('buildSearchParams', () => {
  it('returns an empty object for empty filters', () => {
    expect(buildSearchParams({})).to.deep.equal({});
  });

  it('includes the title when it is set', () => {
    expect(buildSearchParams({ title: 'Le Petit Prince' })).to.deep.equal({
      title: 'Le Petit Prince',
    });
  });

  it('trims spaces around the title', () => {
    expect(buildSearchParams({ title: '  Dune  ' })).to.deep.equal({ title: 'Dune' });
  });

  it('leaves the title out if it only contains spaces', () => {
    expect(buildSearchParams({ title: '   ' })).to.deep.equal({});
  });

  it('includes available=true when the value is true', () => {
    expect(buildSearchParams({ available: true })).to.deep.equal({ available: 'true' });
  });

  it('leaves available out when the value is false', () => {
    expect(buildSearchParams({ available: false })).to.deep.equal({});
  });

  it('keeps barcode as barcode', () => {
    expect(buildSearchParams({ barcode: '0012345' })).to.deep.equal({
      barcode: '0012345',
    });
  });

  it('maps issueNumber to issue_number', () => {
    expect(buildSearchParams({ issueNumber: '42' })).to.deep.equal({ issue_number: '42' });
  });

  it('includes every field set at the same time', () => {
    expect(
      buildSearchParams({
        title: 'Le Monde',
        author: 'Dupont',
        isbn: '978-3-16-148410-0',
        available: true,
      }),
    ).to.deep.equal({
      title: 'Le Monde',
      author: 'Dupont',
      isbn: '978-3-16-148410-0',
      available: 'true',
    });
  });
});

// ─── parseUser ────────────────────────────────────────────────────────────────

describe('parseUser', () => {
  it('returns null for invalid JSON', () => {
    expect(parseUser('not-json')).to.be.null;
  });

  it('returns null for an empty string', () => {
    expect(parseUser('')).to.be.null;
  });

  it('deserializes a valid user', () => {
    const user = { id: 1, last_name: 'Durand', first_name: 'Alice', role: 'user', message: 'ok' };
    expect(parseUser(JSON.stringify(user))).to.deep.equal(user);
  });
});

// ─── isLibrarian ──────────────────────────────────────────────────────────────

describe('isLibrarian', () => {
  it('returns false for null', () => {
    expect(isLibrarian(null)).to.be.false;
  });

  it("returns false for a user with the 'student' role", () => {
    expect(isLibrarian({ id: 1, last_name: 'A', first_name: 'B', role: 'student', message: '' })).to
      .be.false;
  });

  it("returns false for a user with the 'teacher' role", () => {
    expect(isLibrarian({ id: 2, last_name: 'C', first_name: 'D', role: 'teacher', message: '' })).to
      .be.false;
  });

  it("returns true for a user with the 'librarian' role", () => {
    expect(isLibrarian({ id: 3, last_name: 'E', first_name: 'F', role: 'librarian', message: '' }))
      .to.be.true;
  });
});
