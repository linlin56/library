import { inject } from '@angular/core';
import { CanActivateFn, Router } from '@angular/router';
import { ConnexionService } from '../services/connexion.service';

export const connexionGuard: CanActivateFn = () => {
  const connexionService = inject(ConnexionService);
  const router = inject(Router);

  if (connexionService.estConnecte()) {
    return router.createUrlTree(['/']);
  }

  return true;
};
