import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterLink } from '@angular/router';
import { AccueilService, AccueilStats } from '../services/accueil.service';

@Component({
  selector: 'app-accueil',
  standalone: true,
  imports: [CommonModule, RouterLink],
  templateUrl: './accueil.html',
  styleUrl: './accueil.scss'
})
export class AccueilComponent implements OnInit {
  stats: AccueilStats | null = null;
  erreur = false;

  constructor(private accueilService: AccueilService) {}

  ngOnInit(): void {
    this.accueilService.getStats().subscribe({
      next: (data) => (this.stats = data),
      error: () => (this.erreur = true)
    });
  }
}
