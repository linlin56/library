import { Routes } from '@angular/router';
import { AccueilComponent } from './accueil/accueil';

export const routes: Routes = [
  { path: '', component: AccueilComponent },
  { path: '**', redirectTo: '' },
];
