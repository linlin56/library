import { inject } from '@angular/core';
import { CanActivateFn, Router } from '@angular/router';
import { ConnexionService } from '../services/connexion.service';

export const bibliothecaireGuard: CanActivateFn = () => {
  const connexionService = inject(ConnexionService);
  const router = inject(Router);

  const user = connexionService.getUtilisateurConnecte();
  if (user?.role === 'bibliothecaire') {
    return true;
  }
  return router.createUrlTree(['/']);
};
