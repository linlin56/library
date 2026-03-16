import { CommonModule } from '@angular/common';
import { HttpErrorResponse } from '@angular/common/http';
import { Component, inject } from '@angular/core';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { Router } from '@angular/router';
import { finalize } from 'rxjs';
import { ConnexionService } from '../services/connexion.service';

@Component({
  selector: 'app-connexion',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule],
  templateUrl: './connexion.html',
  styleUrl: './connexion.scss',
})
export class ConnexionComponent {
  private readonly formBuilder = inject(FormBuilder);

  erreur = '';
  succes = '';
  soumission = false;

  readonly formulaire = this.formBuilder.nonNullable.group({
    login: ['', [Validators.required]],
    mot_de_passe: ['', [Validators.required]],
  });

  constructor(
    private connexionService: ConnexionService,
    private router: Router,
  ) {}

  soumettre(): void {
    this.erreur = '';
    this.succes = '';

    if (this.formulaire.invalid) {
      this.formulaire.markAllAsTouched();
      this.erreur = 'Veuillez saisir un identifiant et un mot de passe.';
      return;
    }

    this.soumission = true;
    this.connexionService
      .connecter(this.formulaire.getRawValue())
      .pipe(
        finalize(() => {
          this.soumission = false;
        }),
      )
      .subscribe({
        next: (utilisateur) => {
          this.succes = utilisateur.message;
          this.formulaire.reset({ login: '', mot_de_passe: '' });
          void this.router.navigateByUrl('/');
        },
        error: (error: HttpErrorResponse) => {
          this.erreur = error.error?.erreur ?? 'Une erreur est survenue pendant la connexion.';
        },
      });
  }
}
