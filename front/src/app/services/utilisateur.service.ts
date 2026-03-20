import { HttpClient, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';

export interface UtilisateurResume {
  id: number;
  nom: string;
  prenom: string;
  date_de_naissance: string;
  role: string;
}

@Injectable({ providedIn: 'root' })
export class UtilisateurService {
  private readonly apiUrl = 'http://localhost:8080/api/utilisateurs';

  constructor(private http: HttpClient) {}

  rechercherUtilisateurs(
    nom: string,
    prenom: string,
    dateNaissance: string,
  ): Observable<UtilisateurResume[]> {
    let params = new HttpParams();
    if (nom) params = params.set('nom', nom);
    if (prenom) params = params.set('prenom', prenom);
    if (dateNaissance) params = params.set('date_naissance', dateNaissance);
    return this.http.get<UtilisateurResume[]>(`${this.apiUrl}/rechercher`, { params });
  }
}
