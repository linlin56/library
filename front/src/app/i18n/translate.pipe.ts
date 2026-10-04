import { Pipe, PipeTransform, inject } from '@angular/core';
import { LanguageService } from './language.service';

// Impure pipes: they are re-evaluated when the user switches language.

@Pipe({ name: 't', standalone: true, pure: false })
// Translates a key: {{ 'nav.home' | t }} or {{ 'search.resultsCount' | t: { count: 3 } }}
export class TranslatePipe implements PipeTransform {
  private readonly languageService = inject(LanguageService);

  transform(key: string, params?: Record<string, string | number>): string {
    return this.languageService.t(key, params);
  }
}

@Pipe({ name: 'amount', standalone: true, pure: false })
// Formats an amount in euros in the current language: {{ deposit | amount }}
export class AmountPipe implements PipeTransform {
  private readonly languageService = inject(LanguageService);

  transform(amount: number | null | undefined): string {
    return amount == null ? '' : this.languageService.formatAmount(amount);
  }
}
