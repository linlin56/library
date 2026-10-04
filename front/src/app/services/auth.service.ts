import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { BehaviorSubject, Observable, tap } from 'rxjs';

export interface LoginPayload {
  login: string;
  password: string;
}

export interface LoggedInUser {
  id: number;
  last_name: string;
  first_name: string;
  role: string;
  message: string;
}

@Injectable({ providedIn: 'root' })
// Manages the user session: login, logout and persistence in sessionStorage.
export class AuthService {
  private readonly apiUrl = 'http://localhost:8080/api/login';
  private readonly storageKey = 'library.user';
  private readonly loggedInUserSubject = new BehaviorSubject<LoggedInUser | null>(
    this.readUserFromStorage(),
  );

  readonly loggedInUser$ = this.loggedInUserSubject.asObservable();

  constructor(private http: HttpClient) {}

  login(payload: LoginPayload): Observable<LoggedInUser> {
    return this.http.post<LoggedInUser>(this.apiUrl, payload).pipe(
      tap((user) => {
        sessionStorage.setItem(this.storageKey, JSON.stringify(user));
        this.loggedInUserSubject.next(user);
      }),
    );
  }

  getLoggedInUser(): LoggedInUser | null {
    return this.loggedInUserSubject.value;
  }

  isLoggedIn(): boolean {
    return this.getLoggedInUser() !== null;
  }

  logout(): void {
    sessionStorage.removeItem(this.storageKey);
    this.loggedInUserSubject.next(null);
  }

  private readUserFromStorage(): LoggedInUser | null {
    const raw = sessionStorage.getItem(this.storageKey);
    if (!raw) return null;
    try {
      return JSON.parse(raw) as LoggedInUser;
    } catch {
      sessionStorage.removeItem(this.storageKey);
      return null;
    }
  }
}
