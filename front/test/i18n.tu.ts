import { expect } from 'chai';
import { describe, it } from 'mocha';

import {
  DEFAULT_LANGUAGE,
  Dictionary,
  TRANSLATIONS,
  isLanguage,
  translate,
} from '../src/app/i18n/translations';

// Lists every key of a dictionary as 'section.key'
function keys(dictionary: Dictionary, prefix = ''): string[] {
  return Object.entries(dictionary).flatMap(([key, value]) =>
    typeof value === 'string' ? [prefix + key] : keys(value, `${prefix}${key}.`),
  );
}

// ─── Dictionaries ─────────────────────────────────────────────────────────────

describe('TRANSLATIONS', () => {
  it('English is the default language', () => {
    expect(DEFAULT_LANGUAGE).to.equal('en');
  });

  it('French and English define exactly the same keys', () => {
    expect(keys(TRANSLATIONS.en).sort()).to.deep.equal(keys(TRANSLATIONS.fr).sort());
  });

  it('no translation is empty', () => {
    for (const language of ['en', 'fr'] as const) {
      for (const key of keys(TRANSLATIONS[language])) {
        expect(translate(TRANSLATIONS[language], key), `${language}:${key}`).to.not.equal('');
      }
    }
  });
});

// ─── translate ────────────────────────────────────────────────────────────────

describe('translate', () => {
  it('returns the text in the requested language', () => {
    expect(translate(TRANSLATIONS.en, 'nav.home')).to.equal('Home');
    expect(translate(TRANSLATIONS.fr, 'nav.home')).to.equal('Accueil');
  });

  it('replaces the parameters', () => {
    expect(translate(TRANSLATIONS.en, 'search.resultsCount', { count: 3 })).to.equal(
      '3 item(s) found',
    );
  });

  it('leaves a missing parameter untouched', () => {
    expect(translate(TRANSLATIONS.en, 'nav.hello')).to.equal('Hello {name}');
  });

  it('returns the key when it is unknown or incomplete', () => {
    expect(translate(TRANSLATIONS.en, 'unknown.key')).to.equal('unknown.key');
    expect(translate(TRANSLATIONS.en, 'nav')).to.equal('nav');
  });
});

// ─── isLanguage ───────────────────────────────────────────────────────────────

describe('isLanguage', () => {
  it('only accepts en and fr', () => {
    expect(isLanguage('en')).to.be.true;
    expect(isLanguage('fr')).to.be.true;
    expect(isLanguage('de')).to.be.false;
    expect(isLanguage(null)).to.be.false;
  });
});
