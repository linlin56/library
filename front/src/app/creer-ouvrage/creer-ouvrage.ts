import { CommonModule } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { Auteur, CreerOuvrageService, Emplacement } from '../services/creer-ouvrage.service';

@Component({
  selector: 'app-creer-ouvrage',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './creer-ouvrage.html',
  styleUrl: './creer-ouvrage.scss',
})
export class CreerOuvrageComponent implements OnInit {
  type: 'livre' | 'revue' = 'livre';

  // Champs communs
  titre = '';
  caution: number | null = null;

  // Champs livre — auteur autocomplete
  isbn = '';
  auteurQuery = ''; // texte tapé dans le champ
  auteursSuggestions: Auteur[] = []; // suggestions filtrées
  auteurSelectionne: Auteur | null = null; // auteur existant choisi
  suggestionOuverte = false;
  auteurNouveauVisible = false; // champs prenom/nom d'un nouvel auteur
  auteurNom = '';
  auteurPrenom = '';

  // Champs revue
  numero: number | null = null;
  dateParution = '';

  // Etat
  auteurs: Auteur[] = [];
  emplacements: Emplacement[] = [];
  emplacementId: number | null = null;
  enCours = false;
  erreur: string | null = null;
  succes: string | null = null;

  constructor(private service: CreerOuvrageService) {}

  ngOnInit(): void {
    this.service.getAuteurs().subscribe({
      next: (auteurs) => (this.auteurs = auteurs),
    });
    this.service.getEmplacements().subscribe({
      next: (emplacements) => (this.emplacements = emplacements),
    });
  }

  changerType(): void {
    this.erreur = null;
    this.succes = null;
  }

  filtrerAuteurs(): void {
    this.auteurSelectionne = null;
    this.auteurNouveauVisible = false;
    const q = this.auteurQuery.trim().toLowerCase();
    if (!q) {
      this.auteursSuggestions = [];
      this.suggestionOuverte = false;
      return;
    }
    this.auteursSuggestions = this.auteurs.filter(
      (a) => a.nom.toLowerCase().includes(q) || a.prenom.toLowerCase().includes(q),
    );
    this.suggestionOuverte = true;
  }

  selectionnerAuteur(auteur: Auteur): void {
    this.auteurSelectionne = auteur;
    this.auteurQuery = `${auteur.prenom} ${auteur.nom}`;
    this.auteursSuggestions = [];
    this.suggestionOuverte = false;
    this.auteurNouveauVisible = false;
  }

  choisirNouvelAuteur(): void {
    this.auteurSelectionne = null;
    this.auteursSuggestions = [];
    this.suggestionOuverte = false;
    this.auteurNouveauVisible = true;
    this.auteurNom = '';
    this.auteurPrenom = '';
  }

  fermerSuggestions(): void {
    // délai pour laisser le clic sur une suggestion se déclencher avant
    setTimeout(() => (this.suggestionOuverte = false), 150);
  }

  soumettre(): void {
    this.erreur = null;
    this.succes = null;

    if (!this.titre.trim()) {
      this.erreur = 'Le titre est requis.';
      return;
    }
    if (this.caution === null || this.caution < 0) {
      this.erreur = 'La caution doit être un montant positif ou nul.';
      return;
    }

    this.enCours = true;

    if (this.type === 'revue') {
      if (!this.numero || this.numero <= 0) {
        this.erreur = 'Le numéro de la revue est requis.';
        this.enCours = false;
        return;
      }
      if (!this.dateParution) {
        this.erreur = 'La date de parution est requise.';
        this.enCours = false;
        return;
      }
      if (!this.emplacementId) {
        this.erreur = "L'emplacement est requis.";
        this.enCours = false;
        return;
      }
      this.service
        .creerRevue({
          titre: this.titre.trim(),
          caution: this.caution,
          numero: this.numero,
          date_parution: this.dateParution,
          emplacement_id: this.emplacementId,
        })
        .subscribe({
          next: () => {
            this.succes = `La revue « ${this.titre.trim()} » a été créée avec succès.`;
            this.enCours = false;
            this.resetChamps();
          },
          error: (err) => {
            this.erreur = err.error?.erreur ?? 'Erreur lors de la création de la revue.';
            this.enCours = false;
          },
        });
    } else {
      if (!this.isbn.trim()) {
        this.erreur = "L'ISBN est requis.";
        this.enCours = false;
        return;
      }
      if (!this.emplacementId) {
        this.erreur = "L'emplacement est requis.";
        this.enCours = false;
        return;
      }
      const payload: Parameters<CreerOuvrageService['creerLivre']>[0] = {
        titre: this.titre.trim(),
        caution: this.caution,
        isbn: this.isbn.trim(),
        emplacement_id: this.emplacementId,
      };
      if (this.auteurSelectionne) {
        payload.auteur_id = this.auteurSelectionne.id;
      } else if (this.auteurNouveauVisible) {
        if (!this.auteurNom.trim() || !this.auteurPrenom.trim()) {
          this.erreur = "Le nom et le prénom de l'auteur sont requis.";
          this.enCours = false;
          return;
        }
        payload.auteur_nom = this.auteurNom.trim();
        payload.auteur_prenom = this.auteurPrenom.trim();
      } else {
        this.erreur = 'Veuillez choisir ou créer un auteur.';
        this.enCours = false;
        return;
      }
      this.service.creerLivre(payload).subscribe({
        next: () => {
          this.succes = `Le livre « ${this.titre.trim()} » a été créé avec succès.`;
          this.enCours = false;
          const etaitNouvelAuteur = this.auteurNouveauVisible;
          this.resetChamps();
          if (etaitNouvelAuteur) {
            this.service.getAuteurs().subscribe({ next: (a) => (this.auteurs = a) });
          }
        },
        error: (err) => {
          this.erreur = err.error?.erreur ?? 'Erreur lors de la création du livre.';
          this.enCours = false;
        },
      });
    }
  }

  private resetChamps(): void {
    this.titre = '';
    this.caution = null;
    this.isbn = '';
    this.auteurQuery = '';
    this.auteurSelectionne = null;
    this.auteursSuggestions = [];
    this.suggestionOuverte = false;
    this.auteurNouveauVisible = false;
    this.auteurNom = '';
    this.auteurPrenom = '';
    this.numero = null;
    this.dateParution = '';
    this.emplacementId = null;
  }
}
