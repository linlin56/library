import { CommonModule } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { ConnexionService } from '../services/connexion.service';
import { EmpruntItem, EmpruntService } from '../services/emprunt.service';
import { CautionInfo, UtilisateurService } from '../services/utilisateur.service';

@Component({
  selector: 'app-mes-emprunts',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './mes-emprunts.html',
  styleUrl: './mes-emprunts.scss',
})
// Affiche les emprunts actifs de l'utilisateur connecté.
export class MesEmpruntsComponent implements OnInit {
  emprunts: EmpruntItem[] = [];
  enCours = true;
  erreur: string | null = null;
  nomUtilisateur = '';
  caution: CautionInfo | null = null;
  renduEnCours: Set<number> = new Set();

  constructor(
    private empruntService: EmpruntService,
    private connexionService: ConnexionService,
    private utilisateurService: UtilisateurService,
  ) {}

  ngOnInit(): void {
    const u = this.connexionService.getUtilisateurConnecte();
    if (!u) return;
    this.nomUtilisateur = `${u.prenom} ${u.nom}`;
    this.empruntService.listerEmprunts(u.id).subscribe({
      next: (items) => {
        this.emprunts = items;
        this.enCours = false;
      },
      error: () => {
        this.erreur = 'Impossible de charger vos emprunts.';
        this.enCours = false;
      },
    });
    this.utilisateurService.getCaution(u.id).subscribe({
      next: (info) => (this.caution = info),
    });
  }

  isEnRetard(emprunt: EmpruntItem): boolean {
    return emprunt.en_retard;
  }

  declarerRendu(emprunt: EmpruntItem): void {
    this.renduEnCours.add(emprunt.id);
    this.empruntService.rendreLivre(emprunt.id).subscribe({
      next: () => {
        this.emprunts = this.emprunts.filter((e) => e.id !== emprunt.id);
        this.renduEnCours.delete(emprunt.id);
        const u = this.connexionService.getUtilisateurConnecte();
        if (u) {
          this.utilisateurService.getCaution(u.id).subscribe({
            next: (info) => (this.caution = info),
          });
        }
      },
      error: () => {
        this.renduEnCours.delete(emprunt.id);
      },
    });
  }
}
