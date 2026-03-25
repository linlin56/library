import { bootstrapApplication } from '@angular/platform-browser';
import { App } from './app/app';
import { appConfig } from './app/app.config';

// Le fichier main.ts est le point d'entrée de l'application Angular.
// Il utilise la fonction bootstrapApplication pour démarrer l'application en chargeant le composant racine App
// et en appliquant la configuration définie dans appConfig.
bootstrapApplication(App, appConfig).catch((err) => console.error(err));
