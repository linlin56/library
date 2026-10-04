import { DOCUMENT } from '@angular/common';
import { Injectable, inject, signal } from '@angular/core';
import { Title } from '@angular/platform-browser';
import { DEFAULT_LANGUAGE, Language, TRANSLATIONS, isLanguage, translate } from './translations';

@Injectable({ providedIn: 'root' })
// Manages the UI language (English by default) and remembers it in localStorage.
export class LanguageService {
  private readonly storageKey = 'library.language';
  private readonly document = inject(DOCUMENT);
  private readonly title = inject(Title);

  readonly language = signal<Language>(this.readStoredLanguage());

  constructor() {
    this.apply(this.language());
  }

  setLanguage(language: Language): void {
    this.language.set(language);
    try {
      localStorage.setItem(this.storageKey, language);
    } catch {
      // storage unavailable (private browsing...): the language still applies to the session
    }
    this.apply(language);
  }

  t(key: string, params?: Record<string, string | number>): string {
    return translate(TRANSLATIONS[this.language()], key, params);
  }

  /** Translated label of a role returned by the API (student, librarian...). */
  roleLabel(role: string): string {
    const key = `roles.${role}`;
    const label = this.t(key);
    return label === key ? role : label;
  }

  /** Formats an amount in euros for the current language (en: €5.00, fr: 5,00 €). */
  formatAmount(amount: number): string {
    const locale = this.language() === 'fr' ? 'fr-FR' : 'en-GB';
    return new Intl.NumberFormat(locale, { style: 'currency', currency: 'EUR' }).format(amount);
  }

  private apply(language: Language): void {
    this.document.documentElement.lang = language;
    this.title.setTitle(this.t('nav.brand'));
  }

  private readStoredLanguage(): Language {
    try {
      const value = localStorage.getItem(this.storageKey);
      return isLanguage(value) ? value : DEFAULT_LANGUAGE;
    } catch {
      return DEFAULT_LANGUAGE;
    }
  }
}
