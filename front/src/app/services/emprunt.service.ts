import { HttpClient, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';

export interface PreviewEmprunt {
  titre: string;
  code_barre: string;
  caution: number;
  solde_actuel: number;
  nouveau_solde: number;
}

@Injectable({ providedIn: 'root' })
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
}
