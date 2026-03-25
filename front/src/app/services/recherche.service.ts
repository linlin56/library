import { HttpClient, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';

export interface Ouvrage {
  id: number;
  titre: string;
  caution: number;
  type: string; // 'livre' | 'revue'
  isbn?: string;
  auteur?: string;
  numero?: number;
  exemplaires_disponibles: number;
  emplacement?: string;
}

export interface FiltresRecherche {
  titre?: string;
  auteur?: string;
  isbn?: string;
  codeBarre?: string;
  codeRevue?: string;
  disponible?: boolean;
}

@Injectable({ providedIn: 'root' })
// Recherche des ouvrages en appliquant les filtres fournis en paramètres HTTP.
export class RechercheService {
  private readonly apiUrl = 'http://localhost:8080/api/ouvrages';

  constructor(private http: HttpClient) {}

  rechercher(filtres: FiltresRecherche = {}): Observable<Ouvrage[]> {
    let params = new HttpParams();
    if (filtres.titre?.trim()) params = params.set('titre', filtres.titre.trim());
    if (filtres.auteur?.trim()) params = params.set('auteur', filtres.auteur.trim());
    if (filtres.isbn?.trim()) params = params.set('isbn', filtres.isbn.trim());
    if (filtres.codeBarre?.trim()) params = params.set('code_barre', filtres.codeBarre.trim());
    if (filtres.codeRevue?.trim()) params = params.set('code_revue', filtres.codeRevue.trim());
    if (filtres.disponible) params = params.set('disponible', 'true');
    return this.http.get<Ouvrage[]>(this.apiUrl, { params });
  }
}
