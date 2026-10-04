import { HttpClient, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';

export interface UserSummary {
  id: number;
  last_name: string;
  first_name: string;
  phone_number: string;
  role: string;
}

export interface DepositInfo {
  deposit_balance: number;
  total_deposit: number;
}

@Injectable({ providedIn: 'root' })
// Searches users by criteria (librarian view) and manages their deposit.
export class UserService {
  private readonly apiUrl = 'http://localhost:8080/api/users';

  constructor(private http: HttpClient) {}

  searchUsers(
    lastName: string,
    firstName: string,
    postalCode: string,
    phoneNumber: string,
  ): Observable<UserSummary[]> {
    let params = new HttpParams();
    if (lastName) params = params.set('last_name', lastName);
    if (firstName) params = params.set('first_name', firstName);
    if (postalCode) params = params.set('postal_code', postalCode);
    if (phoneNumber) params = params.set('phone_number', phoneNumber);
    return this.http
      .get<UserSummary[]>(`${this.apiUrl}/search`, { params })
      .pipe(map((res) => res ?? []));
  }

  getDeposit(userId: number): Observable<DepositInfo> {
    return this.http.get<DepositInfo>(`${this.apiUrl}/${userId}/deposit`);
  }

  updateTotalDeposit(userId: number, totalDeposit: number): Observable<DepositInfo> {
    return this.http.put<DepositInfo>(`${this.apiUrl}/${userId}/deposit`, {
      total_deposit: totalDeposit,
    });
  }
}
