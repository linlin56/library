import { CommonModule } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { EmpruntEnRetardItem, EmpruntService } from '../services/emprunt.service';

@Component({
  selector: 'app-emprunts-retard',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './emprunts-retard.html',
  styleUrl: './emprunts-retard.scss',
})
export class EmpruntsRetardComponent implements OnInit {
  emprunts: EmpruntEnRetardItem[] = [];
  enCours = true;
  erreur: string | null = null;
  notification: string | null = null;

  constructor(private empruntService: EmpruntService) {}

  ngOnInit(): void {
    this.empruntService.listerEmpruntsEnRetard().subscribe({
      next: (items) => {
        this.emprunts = items;
        this.enCours = false;
      },
      error: () => {
        this.erreur = 'Impossible de charger les emprunts en retard.';
        this.enCours = false;
      },
    });
  }

  envoyerRappels(): void {
    this.notification = 'Mails de rappels envoyés !';
    setTimeout(() => (this.notification = null), 4000);
  }
}
