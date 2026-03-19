import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';

export interface Ouvrage {
  id: number;
  titre: string;
  caution: number;
  isbn?: string;
  auteur_id?: number;
  numero?: number;
}

@Injectable({ providedIn: 'root' })
export class RechercheService {
  private readonly apiUrl = 'http://localhost:8080/api/ouvrages';

  constructor(private http: HttpClient) {}

  rechercher(titre?: string): Observable<Ouvrage[]> {
    let url = this.apiUrl;
    if (titre && titre.trim()) {
      url += `?titre=${encodeURIComponent(titre.trim())}`;
    }
    return this.http.get<Ouvrage[]>(url);
  }
}
