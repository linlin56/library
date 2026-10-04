import { HttpClient, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';

export interface ItemLocation {
  bay_number: number;
  shelf_number: number;
  level: number;
}

export interface Item {
  id: number;
  title: string;
  deposit: number;
  type: string; // 'book' | 'journal'
  isbn?: string;
  author?: string;
  issue_number?: number;
  available_copies: number;
  location?: ItemLocation;
}

export interface SearchFilters {
  title?: string;
  author?: string;
  isbn?: string;
  barcode?: string;
  issueNumber?: string;
  available?: boolean;
}

@Injectable({ providedIn: 'root' })
// Searches items by sending the given filters as HTTP parameters.
export class SearchService {
  private readonly apiUrl = 'http://localhost:8080/api/items';

  constructor(private http: HttpClient) {}

  search(filters: SearchFilters = {}): Observable<Item[]> {
    let params = new HttpParams();
    if (filters.title?.trim()) params = params.set('title', filters.title.trim());
    if (filters.author?.trim()) params = params.set('author', filters.author.trim());
    if (filters.isbn?.trim()) params = params.set('isbn', filters.isbn.trim());
    if (filters.barcode?.trim()) params = params.set('barcode', filters.barcode.trim());
    if (filters.issueNumber?.trim())
      params = params.set('issue_number', filters.issueNumber.trim());
    if (filters.available) params = params.set('available', 'true');
    return this.http.get<Item[]>(this.apiUrl, { params }).pipe(map((res) => res ?? []));
  }
}
