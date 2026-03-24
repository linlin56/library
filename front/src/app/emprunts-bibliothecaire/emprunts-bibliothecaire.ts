import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { EmpruntItem, EmpruntService } from '../services/emprunt.service';
import { UtilisateurResume, UtilisateurService } from '../services/utilisateur.service';

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
