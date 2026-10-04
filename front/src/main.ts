import { bootstrapApplication } from '@angular/platform-browser';
import { App } from './app/app';
import { appConfig } from './app/app.config';

// main.ts is the entry point of the Angular application.
// It uses bootstrapApplication to start the application by loading the root App component
// and applying the configuration defined in appConfig.
bootstrapApplication(App, appConfig).catch((err) => console.error(err));
