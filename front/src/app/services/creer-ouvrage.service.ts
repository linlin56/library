import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';

export interface Auteur {
  id: number;
  nom: string;
  prenom: string;
}

export interface CreerOuvrageResponse {
  id: number;
  message: string;
}

@Injectable({ providedIn: 'root' })
export class CreerOuvrageService {
  private readonly baseUrl = 'http://localhost:8080/api';

  constructor(private http: HttpClient) {}

  getAuteurs(): Observable<Auteur[]> {
    return this.http.get<Auteur[]>(`${this.baseUrl}/auteurs`);
  }

  creerLivre(data: {
    titre: string;
    caution: number;
    isbn: string;
    auteur_id?: number;
    auteur_nom?: string;
    auteur_prenom?: string;
  }): Observable<CreerOuvrageResponse> {
    return this.http.post<CreerOuvrageResponse>(`${this.baseUrl}/livres`, data);
  }

  creerRevue(data: {
    titre: string;
    caution: number;
    numero: number;
  }): Observable<CreerOuvrageResponse> {
    return this.http.post<CreerOuvrageResponse>(`${this.baseUrl}/revues`, data);
  }
}
