import { inject } from '@angular/core';
import { CanActivateFn, Router } from '@angular/router';
import { ConnexionService } from '../services/connexion.service';

// Vérifie qu'un utilisateur est connecté. Redirige vers /connexion sinon.
export const authGuard: CanActivateFn = () => {
  const connexionService = inject(ConnexionService);
  const router = inject(Router);

  if (connexionService.estConnecte()) {
    return true;
  }
  return router.createUrlTree(['/connexion']);
};
