import { CommonModule } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { RouterLink } from '@angular/router';
import { AccueilService, AccueilStats } from '../services/accueil.service';
import { ConnexionService } from '../services/connexion.service';

@Component({
  selector: 'app-accueil',
  standalone: true,
  imports: [CommonModule, RouterLink],
  templateUrl: './accueil.html',
  styleUrl: './accueil.scss',
})
export class AccueilComponent implements OnInit {
  stats: AccueilStats | null = null;
  erreur = false;
  estConnecte = false;

  constructor(
    private accueilService: AccueilService,
    private connexionService: ConnexionService,
  ) {}

  ngOnInit(): void {
    this.estConnecte = this.connexionService.estConnecte();
    this.accueilService.getStats().subscribe({
      next: (data) => (this.stats = data),
      error: () => (this.erreur = true),
    });
  }
}
