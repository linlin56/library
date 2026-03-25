import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { BehaviorSubject, Observable, tap } from 'rxjs';

export interface ConnexionPayload {
  login: string;
  mot_de_passe: string;
}

export interface UtilisateurConnecte {
  id: number;
  nom: string;
  prenom: string;
  role: string;
  message: string;
}

@Injectable({ providedIn: 'root' })
// Gère la session utilisateur : connexion, déconnexion et persistance via sessionStorage.
export class ConnexionService {
  private readonly apiUrl = 'http://localhost:8080/api/connexion';
  private readonly storageKey = 'bibliotheque.utilisateur';
  private readonly utilisateurConnecteSubject = new BehaviorSubject<UtilisateurConnecte | null>(
    this.readUtilisateurFromStorage(),
  );

  readonly utilisateurConnecte$ = this.utilisateurConnecteSubject.asObservable();

  constructor(private http: HttpClient) {}

  connecter(payload: ConnexionPayload): Observable<UtilisateurConnecte> {
    return this.http.post<UtilisateurConnecte>(this.apiUrl, payload).pipe(
      tap((utilisateur) => {
        sessionStorage.setItem(this.storageKey, JSON.stringify(utilisateur));
        this.utilisateurConnecteSubject.next(utilisateur);
      }),
    );
  }

  getUtilisateurConnecte(): UtilisateurConnecte | null {
    return this.utilisateurConnecteSubject.value;
  }

  estConnecte(): boolean {
    return this.getUtilisateurConnecte() !== null;
  }

  deconnecter(): void {
    sessionStorage.removeItem(this.storageKey);
    this.utilisateurConnecteSubject.next(null);
  }

  private readUtilisateurFromStorage(): UtilisateurConnecte | null {
    const raw = sessionStorage.getItem(this.storageKey);
    if (!raw) return null;
    try {
      return JSON.parse(raw) as UtilisateurConnecte;
    } catch {
      sessionStorage.removeItem(this.storageKey);
      return null;
    }
  }
}
