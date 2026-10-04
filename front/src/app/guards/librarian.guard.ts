import { inject } from '@angular/core';
import { CanActivateFn, Router } from '@angular/router';
import { AuthService } from '../services/auth.service';

// Only lets logged-in librarians through. Redirects to the home page otherwise.
export const librarianGuard: CanActivateFn = () => {
  const authService = inject(AuthService);
  const router = inject(Router);

  const user = authService.getLoggedInUser();
  if (user?.role === 'librarian') {
    return true;
  }
  return router.createUrlTree(['/']);
};
