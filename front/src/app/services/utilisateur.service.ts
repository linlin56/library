import { HttpClient, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';

export interface UtilisateurResume {
  id: number;
  nom: string;
  prenom: string;
  numero_telephone: string;
  role: string;
}

export interface CautionInfo {
  solde_caution: number;
  caution_totale: number;
}

@Injectable({ providedIn: 'root' })
// Recherche des utilisateurs par critères (vue bibliothécaire).
export class UtilisateurService {
  private readonly apiUrl = 'http://localhost:8080/api/utilisateurs';

  constructor(private http: HttpClient) {}

  rechercherUtilisateurs(
    nom: string,
    prenom: string,
    codePostal: string,
    numeroTelephone: string,
  ): Observable<UtilisateurResume[]> {
    let params = new HttpParams();
    if (nom) params = params.set('nom', nom);
    if (prenom) params = params.set('prenom', prenom);
    if (codePostal) params = params.set('code_postal', codePostal);
    if (numeroTelephone) params = params.set('numero_telephone', numeroTelephone);
    return this.http.get<UtilisateurResume[]>(`${this.apiUrl}/rechercher`, { params });
  }

  getCaution(utilisateurId: number): Observable<CautionInfo> {
    return this.http.get<CautionInfo>(`${this.apiUrl}/${utilisateurId}/caution`);
  }

  updateCautionTotale(utilisateurId: number, cautionTotale: number): Observable<CautionInfo> {
    return this.http.put<CautionInfo>(`${this.apiUrl}/${utilisateurId}/caution`, {
      caution_totale: cautionTotale,
    });
  }
}
