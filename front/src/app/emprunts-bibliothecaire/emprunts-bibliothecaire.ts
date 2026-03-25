import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { EmpruntItem, EmpruntService } from '../services/emprunt.service';
import { CautionInfo, UtilisateurResume, UtilisateurService } from '../services/utilisateur.service';

type Etape = 'recherche' | 'selection' | 'emprunts';

@Component({
  selector: 'app-emprunts-bibliothecaire',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './emprunts-bibliothecaire.html',
  styleUrl: './emprunts-bibliothecaire.scss',
})
export class EmpruntsBibliothecaireComponent {
  // Formulaire de recherche
  nom = '';
  prenom = '';
  codePostal = '';
  numeroTelephone = '';

  // Etat
  etape: Etape = 'recherche';
  enCours = false;
  erreur: string | null = null;

  // Resultats
  utilisateurs: UtilisateurResume[] = [];
  utilisateurSelectionne: UtilisateurResume | null = null;
  emprunts: EmpruntItem[] = [];
  caution: CautionInfo | null = null;

  // Edition caution
  editionCaution = false;
  nouvelleCautionTotale: number | null = null;
  enregistrementCaution = false;
  erreurCaution: string | null = null;
  successCaution = false;

  constructor(
    private utilisateurService: UtilisateurService,
    private empruntService: EmpruntService,
  ) {}

  peutRechercher(): boolean {
    return !!(
      this.nom.trim() ||
      this.prenom.trim() ||
      this.codePostal.trim() ||
      this.numeroTelephone.trim()
    );
  }

  rechercherUtilisateur(): void {
    if (!this.peutRechercher()) {
      this.erreur = 'Renseignez au moins un critère de recherche.';
      return;
    }

    this.enCours = true;
    this.erreur = null;
    this.utilisateurs = [];
    this.utilisateurSelectionne = null;

    this.utilisateurService
      .rechercherUtilisateurs(
        this.nom.trim(),
        this.prenom.trim(),
        this.codePostal.trim(),
        this.numeroTelephone.trim(),
      )
      .subscribe({
        next: (utilisateurs) => {
          this.utilisateurs = utilisateurs;
          this.enCours = false;
          if (utilisateurs.length === 0) {
            this.erreur = 'Aucun utilisateur trouvé avec ces critères.';
          } else {
            this.etape = 'selection';
          }
        },
        error: () => {
          this.erreur = 'Erreur lors de la recherche.';
          this.enCours = false;
        },
      });
  }

  selectionnerUtilisateur(utilisateur: UtilisateurResume): void {
    this.utilisateurSelectionne = utilisateur;
    this.enCours = true;
    this.emprunts = [];
    this.caution = null;
    this.editionCaution = false;
    this.erreurCaution = null;
    this.successCaution = false;

    this.empruntService.listerEmprunts(utilisateur.id).subscribe({
      next: (emprunts) => {
        this.emprunts = emprunts;
        this.etape = 'emprunts';
        this.enCours = false;
      },
      error: () => {
        this.erreur = 'Impossible de charger les emprunts.';
        this.enCours = false;
      },
    });
    this.utilisateurService.getCaution(utilisateur.id).subscribe({
      next: (info) => (this.caution = info),
    });
  }

  ouvrirEditionCaution(): void {
    this.nouvelleCautionTotale = this.caution?.caution_totale ?? null;
    this.editionCaution = true;
    this.erreurCaution = null;
    this.successCaution = false;
  }

  annulerEditionCaution(): void {
    this.editionCaution = false;
    this.erreurCaution = null;
  }

  enregistrerCaution(): void {
    if (!this.utilisateurSelectionne || this.nouvelleCautionTotale === null || this.nouvelleCautionTotale < 0) {
      this.erreurCaution = 'Valeur invalide.';
      return;
    }
    this.enregistrementCaution = true;
    this.erreurCaution = null;
    this.utilisateurService.updateCautionTotale(this.utilisateurSelectionne.id, this.nouvelleCautionTotale).subscribe({
      next: () => {
        if (this.caution) this.caution = { ...this.caution, caution_totale: this.nouvelleCautionTotale! };
        this.editionCaution = false;
        this.successCaution = true;
        this.enregistrementCaution = false;
      },
      error: () => {
        this.erreurCaution = 'Impossible de mettre à jour la caution.';
        this.enregistrementCaution = false;
      },
    });
  }

  retourRecherche(): void {
    this.etape = 'recherche';
    this.utilisateurSelectionne = null;
    this.emprunts = [];
    this.utilisateurs = [];
    this.erreur = null;
    this.nom = '';
    this.prenom = '';
    this.codePostal = '';
    this.numeroTelephone = '';
  }

  formatRole(role: string): string {
    const roles: Record<string, string> = {
      etudiant: 'Etudiant',
      enseignant: 'Enseignant',
      utilisateur: 'Utilisateur',
    };
    return roles[role] ?? role;
  }
}
