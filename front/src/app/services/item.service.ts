import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';

export interface Author {
  id: number;
  last_name: string;
  first_name: string;
}

export interface Location {
  id: number;
  bay_number: number;
  shelf_number: number;
  level: number;
  category_name: string;
}

export interface CreateItemResponse {
  id: number;
  message: string;
}

@Injectable({ providedIn: 'root' })
// API access to manage items: loading reference data and creating books/journals.
export class ItemService {
  private readonly baseUrl = 'http://localhost:8080/api';

  constructor(private http: HttpClient) {}

  getAuthors(): Observable<Author[]> {
    return this.http.get<Author[]>(`${this.baseUrl}/authors`);
  }

  getLocations(): Observable<Location[]> {
    return this.http.get<Location[]>(`${this.baseUrl}/locations`);
  }

  createBook(data: {
    title: string;
    deposit: number;
    isbn: string;
    location_id: number;
    author_id?: number;
    author_last_name?: string;
    author_first_name?: string;
  }): Observable<CreateItemResponse> {
    return this.http.post<CreateItemResponse>(`${this.baseUrl}/books`, data);
  }

  createJournal(data: {
    title: string;
    deposit: number;
    issue_number: number;
    publication_date: string;
    location_id: number;
  }): Observable<CreateItemResponse> {
    return this.http.post<CreateItemResponse>(`${this.baseUrl}/journals`, data);
  }
}
