import { CommonModule } from '@angular/common';
import { HttpErrorResponse } from '@angular/common/http';
import { Component, inject } from '@angular/core';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { Router } from '@angular/router';
import { finalize } from 'rxjs';
import { LanguageService } from '../i18n/language.service';
import { TranslatePipe } from '../i18n/translate.pipe';
import { AuthService } from '../services/auth.service';

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, TranslatePipe],
  templateUrl: './login.html',
  styleUrl: './login.scss',
})
export class LoginComponent {
  private readonly formBuilder = inject(FormBuilder);
  private readonly languageService = inject(LanguageService);

  error = '';
  success = '';
  submitting = false;

  readonly form = this.formBuilder.nonNullable.group({
    login: ['', [Validators.required]],
    password: ['', [Validators.required]],
  });

  constructor(
    private authService: AuthService,
    private router: Router,
  ) {}

  submit(): void {
    this.error = '';
    this.success = '';

    if (this.form.invalid) {
      this.form.markAllAsTouched();
      this.error = this.languageService.t('login.errors.missing');
      return;
    }

    this.submitting = true;
    this.authService
      .login(this.form.getRawValue())
      .pipe(
        finalize(() => {
          this.submitting = false;
        }),
      )
      .subscribe({
        next: (user) => {
          this.success = user.message;
          this.form.reset({ login: '', password: '' });
          void this.router.navigateByUrl('/');
        },
        error: (error: HttpErrorResponse) => {
          this.error = error.error?.error ?? this.languageService.t('login.errors.generic');
        },
      });
  }
}
