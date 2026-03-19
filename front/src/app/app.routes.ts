import { Routes } from '@angular/router';
import { AccueilComponent } from './accueil/accueil';
import { ConnexionComponent } from './connexion/connexion';
import { CreerUtilisateurComponent } from './creer-utilisateur/creer-utilisateur';
import { bibliothecaireGuard } from './guards/bibliothecaire.guard';
import { connexionGuard } from './guards/connexion.guard';

export const routes: Routes = [
  { path: '', component: AccueilComponent },
  { path: 'connexion', component: ConnexionComponent, canActivate: [connexionGuard] },
  {
    path: 'creer-utilisateur',
    component: CreerUtilisateurComponent,
    canActivate: [bibliothecaireGuard],
  },
  { path: '**', redirectTo: '' },
];
