import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';

export interface HomeStats {
  book_count: number;
  journal_count: number;
  available_copy_count: number;
  user_count: number;
}

@Injectable({ providedIn: 'root' })
// Fetches the statistics displayed on the home page.
export class HomeService {
  private readonly apiUrl = 'http://localhost:8080/api/home';

  constructor(private http: HttpClient) {}

  getStats(): Observable<HomeStats> {
    return this.http.get<HomeStats>(this.apiUrl);
  }
}
