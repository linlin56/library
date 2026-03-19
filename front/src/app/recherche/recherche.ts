import { CommonModule } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { finalize } from 'rxjs/operators';
import { ConnexionService } from '../services/connexion.service';
import { Ouvrage, RechercheService } from '../services/recherche.service';

@Component({
  selector: 'app-recherche',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './recherche.html',
  styleUrl: './recherche.scss',
})
export class RechercheComponent implements OnInit {
  titre = '';
  resultats: Ouvrage[] = [];
  enCours = false;
  erreur: string | null = null;
  aRecherche = false;
  estConnecte = false;

  constructor(
    private rechercheService: RechercheService,
    private connexionService: ConnexionService,
  ) {}

  ngOnInit(): void {
    this.estConnecte = this.connexionService.estConnecte();
  }

  rechercher(): void {
    if (!this.titre.trim()) {
      this.erreur = 'Veuillez entrer un titre à rechercher.';
      return;
    }

    this.enCours = true;
    this.erreur = null;
    this.resultats = [];
    this.aRecherche = false;

    this.rechercheService
      .rechercher(this.titre)
      .pipe(finalize(() => (this.enCours = false)))
      .subscribe({
        next: (ouvrages) => {
          this.resultats = ouvrages ?? [];
          this.aRecherche = true;
        },
        error: () => {
          this.erreur = 'Erreur lors de la recherche. Veuillez réessayer.';
          this.resultats = [];
          this.aRecherche = true;
        },
      });
  }

  effacer(): void {
    this.titre = '';
    this.resultats = [];
    this.erreur = null;
    this.aRecherche = false;
  }

  getTypeBadge(ouvrage: Ouvrage): string {
    return ouvrage.isbn ? 'Livre' : 'Revue';
  }

  getTypeBadgeClass(ouvrage: Ouvrage): string {
    return ouvrage.isbn ? 'bg-info' : 'bg-success';
  }
}
