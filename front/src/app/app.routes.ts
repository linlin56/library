import { Routes } from '@angular/router';
import { AccueilComponent } from './accueil/accueil';
import { ConnexionComponent } from './connexion/connexion';
import { connexionGuard } from './guards/connexion.guard';

export const routes: Routes = [
  { path: '', component: AccueilComponent },
  { path: 'connexion', component: ConnexionComponent, canActivate: [connexionGuard] },
  { path: '**', redirectTo: '' },
];
