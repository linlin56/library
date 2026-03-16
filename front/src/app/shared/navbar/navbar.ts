import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { Router, RouterLink, RouterLinkActive } from '@angular/router';
import { Observable } from 'rxjs';
import { ConnexionService, UtilisateurConnecte } from '../../services/connexion.service';

@Component({
  selector: 'app-navbar',
  standalone: true,
  imports: [CommonModule, RouterLink, RouterLinkActive],
  templateUrl: './navbar.html',
})
export class NavbarComponent {
  readonly utilisateurConnecte$: Observable<UtilisateurConnecte | null>;

  constructor(
    private connexionService: ConnexionService,
    private router: Router,
  ) {
    this.utilisateurConnecte$ = this.connexionService.utilisateurConnecte$;
  }

  deconnecter(): void {
    this.connexionService.deconnecter();
    void this.router.navigateByUrl('/connexion');
  }
}
