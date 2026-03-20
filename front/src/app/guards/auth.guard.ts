import { inject } from '@angular/core';
import { CanActivateFn, Router } from '@angular/router';
import { ConnexionService } from '../services/connexion.service';

export const authGuard: CanActivateFn = () => {
  const connexionService = inject(ConnexionService);
  const router = inject(Router);

  if (connexionService.estConnecte()) {
    return true;
  }
  return router.createUrlTree(['/connexion']);
};
