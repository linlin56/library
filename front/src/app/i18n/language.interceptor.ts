import { HttpInterceptorFn } from '@angular/common/http';
import { inject } from '@angular/core';
import { LanguageService } from './language.service';

// Sends the selected language to the API so that it answers in the same language.
export const languageInterceptor: HttpInterceptorFn = (req, next) => {
  const language = inject(LanguageService).language();
  return next(req.clone({ setHeaders: { 'Accept-Language': language } }));
};
