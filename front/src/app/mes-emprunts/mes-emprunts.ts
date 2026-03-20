import { CommonModule } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { ConnexionService } from '../services/connexion.service';
import { EmpruntItem, EmpruntService } from '../services/emprunt.service';

@Component({
  selector: 'app-mes-emprunts',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './mes-emprunts.html',
  styleUrl: './mes-emprunts.scss',
})
export class MesEmpruntsComponent implements OnInit {
  emprunts: EmpruntItem[] = [];
  enCours = true;
  erreur: string | null = null;
  nomUtilisateur = '';

  constructor(
    private empruntService: EmpruntService,
    private connexionService: ConnexionService,
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
  }

  isEnRetard(emprunt: EmpruntItem): boolean {
    return emprunt.en_retard;
  }
}
