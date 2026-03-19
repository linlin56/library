import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';

export interface Departement {
  Id: number;
  Nom: string;
}

export interface CreerUtilisateurPayload {
  nom: string;
  prenom: string;
  numero_telephone: string;
  date_naissance: string;
  email: string;
  statut: 'etudiant' | 'enseignant';
  annee_etude?: string;
  departement_id?: number;
  numero_rue: string;
  nom_rue: string;
  code_postal: string;
  ville: string;
}

export interface CreerUtilisateurResponse {
  login: string;
  mot_de_passe: string;
  message: string;
}

@Injectable({ providedIn: 'root' })
export class CreerUtilisateurService {
  private readonly apiUrl = 'http://localhost:8080/api';

  constructor(private http: HttpClient) {}

  getDepartements(): Observable<Departement[]> {
    return this.http.get<Departement[]>(`${this.apiUrl}/departements`);
  }

  creerUtilisateur(data: CreerUtilisateurPayload): Observable<CreerUtilisateurResponse> {
    return this.http.post<CreerUtilisateurResponse>(`${this.apiUrl}/utilisateurs`, data);
  }
}
