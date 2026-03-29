import { HttpClient, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';

export interface PreviewEmprunt {
  titre: string;
  code_barre: string;
  caution: number;
  solde_actuel: number;
  nouveau_solde: number;
}

export interface EmpruntItem {
  id: number;
  code_barre: string;
  titre: string;
  date_debut: string;
  date_fin: string;
  en_retard: boolean;
}

export interface ExemplaireDisponible {
  id: number;
  code_barre: string;
}

export interface ExemplaireComplet {
  id: number;
  code_barre: string;
  est_emprunte: boolean;
  date_fin_emprunt?: string;
}

@Injectable({ providedIn: 'root' })
// Gère les emprunts : vérification des droits, enregistrement et consultation.
export class EmpruntService {
  private readonly apiUrl = 'http://localhost:8080/api/emprunts';

  constructor(private http: HttpClient) {}

  verifier(utilisateurId: number, codeBarre: string): Observable<PreviewEmprunt> {
    const params = new HttpParams()
      .set('utilisateur_id', utilisateurId.toString())
      .set('code_barre', codeBarre);
    return this.http.get<PreviewEmprunt>(`${this.apiUrl}/verifier`, { params });
  }

  emprunter(utilisateurId: number, codeBarre: string): Observable<{ message: string }> {
    return this.http.post<{ message: string }>(this.apiUrl, {
      utilisateur_id: utilisateurId,
      code_barre: codeBarre,
    });
  }

  listerEmprunts(utilisateurId: number): Observable<EmpruntItem[]> {
    const params = new HttpParams().set('utilisateur_id', utilisateurId.toString());
    return this.http.get<EmpruntItem[]>(this.apiUrl, { params }).pipe(map((res) => res ?? []));
  }

  listerEmpruntsEnRetard(): Observable<EmpruntEnRetardItem[]> {
    return this.http
      .get<EmpruntEnRetardItem[]>(`${this.apiUrl}/retard`)
      .pipe(map((res) => res ?? []));
  }

  getExemplairesDisponibles(ouvrageId: number): Observable<ExemplaireDisponible[]> {
    return this.http
      .get<ExemplaireDisponible[]>(`http://localhost:8080/api/ouvrages/${ouvrageId}/exemplaires`)
      .pipe(map((res) => res ?? []));
  }

  getTousExemplaires(ouvrageId: number): Observable<ExemplaireComplet[]> {
    return this.http
      .get<ExemplaireComplet[]>(`http://localhost:8080/api/ouvrages/${ouvrageId}/exemplaires/tous`)
      .pipe(map((res) => res ?? []));
  }

  creerExemplaire(ouvrageId: number, codeBarre: string): Observable<{ id: number }> {
    return this.http.post<{ id: number }>(
      `http://localhost:8080/api/ouvrages/${ouvrageId}/exemplaires`,
      { code_barre: codeBarre },
    );
  }
}

export interface EmpruntEnRetardItem {
  id: number;
  code_barre: string;
  titre: string;
  date_fin: string;
  nom: string;
  prenom: string;
  numero_telephone: string;
}
