import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { Router, RouterLink, RouterLinkActive } from '@angular/router';
import { Observable } from 'rxjs';
import { LanguageService } from '../../i18n/language.service';
import { TranslatePipe } from '../../i18n/translate.pipe';
import { Language } from '../../i18n/translations';
import { AuthService, LoggedInUser } from '../../services/auth.service';

@Component({
  selector: 'app-navbar',
  standalone: true,
  imports: [CommonModule, RouterLink, RouterLinkActive, TranslatePipe],
  templateUrl: './navbar.html',
})
export class NavbarComponent {
  readonly loggedInUser$: Observable<LoggedInUser | null>;
  readonly languages: { code: Language; name: string }[] = [
    { code: 'en', name: 'English' },
    { code: 'fr', name: 'Français' },
  ];

  constructor(
    private authService: AuthService,
    private languageService: LanguageService,
    private router: Router,
  ) {
    this.loggedInUser$ = this.authService.loggedInUser$;
  }

  get language() {
    return this.languageService.language;
  }

  setLanguage(language: Language): void {
    this.languageService.setLanguage(language);
  }

  formatRole(role: string): string {
    return this.languageService.roleLabel(role === 'librarian' ? 'librarian' : 'user');
  }

  logout(): void {
    this.authService.logout();
    void this.router.navigateByUrl('/login');
  }
}
