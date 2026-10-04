import { HttpClient, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';

export interface LoanPreview {
  title: string;
  barcode: string;
  deposit: number;
  current_balance: number;
  new_balance: number;
}

export interface Loan {
  id: number;
  barcode: string;
  title: string;
  start_date: string;
  due_date: string;
  overdue: boolean;
}

export interface OverdueLoan {
  id: number;
  barcode: string;
  title: string;
  due_date: string;
  last_name: string;
  first_name: string;
  phone_number: string;
}

export interface AvailableCopy {
  id: number;
  barcode: string;
}

export interface CopyDetail {
  id: number;
  barcode: string;
  is_borrowed: boolean;
  loan_due_date?: string;
}

@Injectable({ providedIn: 'root' })
// Manages loans: eligibility check, recording and listing, plus copies.
export class LoanService {
  private readonly apiUrl = 'http://localhost:8080/api/loans';
  private readonly itemsUrl = 'http://localhost:8080/api/items';

  constructor(private http: HttpClient) {}

  check(userId: number, barcode: string): Observable<LoanPreview> {
    const params = new HttpParams().set('user_id', userId.toString()).set('barcode', barcode);
    return this.http.get<LoanPreview>(`${this.apiUrl}/check`, { params });
  }

  borrow(userId: number, barcode: string): Observable<{ message: string }> {
    return this.http.post<{ message: string }>(this.apiUrl, {
      user_id: userId,
      barcode: barcode,
    });
  }

  listLoans(userId: number): Observable<Loan[]> {
    const params = new HttpParams().set('user_id', userId.toString());
    return this.http.get<Loan[]>(this.apiUrl, { params }).pipe(map((res) => res ?? []));
  }

  listOverdueLoans(): Observable<OverdueLoan[]> {
    return this.http.get<OverdueLoan[]>(`${this.apiUrl}/overdue`).pipe(map((res) => res ?? []));
  }

  getAvailableCopies(itemId: number): Observable<AvailableCopy[]> {
    return this.http
      .get<AvailableCopy[]>(`${this.itemsUrl}/${itemId}/copies`)
      .pipe(map((res) => res ?? []));
  }

  getAllCopies(itemId: number): Observable<CopyDetail[]> {
    return this.http
      .get<CopyDetail[]>(`${this.itemsUrl}/${itemId}/copies/all`)
      .pipe(map((res) => res ?? []));
  }

  createCopy(itemId: number, barcode: string): Observable<{ id: number }> {
    return this.http.post<{ id: number }>(`${this.itemsUrl}/${itemId}/copies`, { barcode });
  }

  returnCopy(copyId: number): Observable<{ message: string }> {
    return this.http.delete<{ message: string }>(`${this.apiUrl}/${copyId}`);
  }
}
