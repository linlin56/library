import { CommonModule } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { ConnexionService } from '../services/connexion.service';
import { FiltresRecherche, Ouvrage, RechercheService } from '../services/recherche.service';

@Component({
  selector: 'app-recherche',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './recherche.html',
  styleUrl: './recherche.scss',
})
export class RechercheComponent implements OnInit {
  // Champs de recherche
  titre = '';
  auteur = '';
  isbn = '';
  codeRevue = '';
  codeBarre = '';
  disponible = false;
  rechercheAvanceeOuverte = false;

  // Etat
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
    const auMoinsUnChamp =
      this.titre.trim() ||
      this.auteur.trim() ||
      this.isbn.trim() ||
      this.codeRevue.trim() ||
      this.codeBarre.trim();

    if (!auMoinsUnChamp) {
      this.erreur = 'Veuillez renseigner au moins un champ.';
      return;
    }

    this.enCours = true;
    this.erreur = null;
    this.resultats = [];
    this.aRecherche = false;

    const filtres: FiltresRecherche = {
      titre: this.titre,
      auteur: this.auteur,
      isbn: this.isbn,
      codeRevue: this.codeRevue,
      codeBarre: this.codeBarre,
      disponible: this.disponible,
    };

    this.rechercheService.rechercher(filtres).subscribe({
      next: (ouvrages) => {
        this.resultats = ouvrages ?? [];
        this.aRecherche = true;
        this.enCours = false;
      },
      error: () => {
        this.erreur = 'Erreur lors de la recherche. Veuillez reessayer.';
        this.resultats = [];
        this.aRecherche = true;
        this.enCours = false;
      },
    });
  }

  effacer(): void {
    this.titre = '';
    this.auteur = '';
    this.isbn = '';
    this.codeRevue = '';
    this.codeBarre = '';
    this.disponible = false;
    this.resultats = [];
    this.erreur = null;
    this.aRecherche = false;
    this.enCours = false;
  }

  getTypeBadge(ouvrage: Ouvrage): string {
    return ouvrage.type === 'livre' ? 'Livre' : 'Revue';
  }

  getTypeBadgeClass(ouvrage: Ouvrage): string {
    return ouvrage.type === 'livre' ? 'bg-info' : 'bg-success';
  }
}
