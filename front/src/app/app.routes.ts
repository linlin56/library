import { Routes } from '@angular/router';
import { AccueilComponent } from './accueil/accueil';
import { ConnexionComponent } from './connexion/connexion';
import { CreerOuvrageComponent } from './creer-ouvrage/creer-ouvrage';
import { CreerUtilisateurComponent } from './creer-utilisateur/creer-utilisateur';
import { EmpruntsBibliothecaireComponent } from './emprunts-bibliothecaire/emprunts-bibliothecaire';
import { EmpruntsRetardComponent } from './emprunts-retard/emprunts-retard';
import { authGuard } from './guards/auth.guard';
import { bibliothecaireGuard } from './guards/bibliothecaire.guard';
import { connexionGuard } from './guards/connexion.guard';
import { MesEmpruntsComponent } from './mes-emprunts/mes-emprunts';
import { RechercheComponent } from './recherche/recherche';

export const routes: Routes = [
  { path: '', component: AccueilComponent },
  { path: 'connexion', component: ConnexionComponent, canActivate: [connexionGuard] },
  {
    path: 'creer-utilisateur',
    component: CreerUtilisateurComponent,
    canActivate: [bibliothecaireGuard],
  },
  {
    path: 'creer-ouvrage',
    component: CreerOuvrageComponent,
    canActivate: [bibliothecaireGuard],
  },
  { path: 'recherche', component: RechercheComponent },
  { path: 'mes-emprunts', component: MesEmpruntsComponent, canActivate: [authGuard] },
  {
    path: 'emprunts-bibliothecaire',
    component: EmpruntsBibliothecaireComponent,
    canActivate: [bibliothecaireGuard],
  },
  {
    path: 'emprunts-retard',
    component: EmpruntsRetardComponent,
    canActivate: [bibliothecaireGuard],
  },
  { path: '**', redirectTo: '' },
];
