import { inject } from '@angular/core';
import { CanActivateFn, Router } from '@angular/router';
import { ConnexionService } from '../services/connexion.service';

// Empêche un utilisateur déjà connecté d'accéder à la page de connexion.
export const connexionGuard: CanActivateFn = () => {
  const connexionService = inject(ConnexionService);
  const router = inject(Router);

  if (connexionService.estConnecte()) {
    return router.createUrlTree(['/']);
  }

  return true;
};
