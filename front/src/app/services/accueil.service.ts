import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

export interface AccueilStats {
  nb_livres: number;
  nb_revues: number;
  nb_exemplaires_disponibles: number;
  nb_utilisateurs: number;
}

@Injectable({ providedIn: 'root' })
export class AccueilService {
  private readonly apiUrl = 'http://localhost:8080/api/accueil';

  constructor(private http: HttpClient) {}

  getStats(): Observable<AccueilStats> {
    return this.http.get<AccueilStats>(this.apiUrl);
  }
}
