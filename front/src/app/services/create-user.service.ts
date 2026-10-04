import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';

export interface Department {
  id: number;
  name: string;
}

export interface CreateUserPayload {
  last_name: string;
  first_name: string;
  phone_number: string;
  birth_date: string;
  email: string;
  status: 'student' | 'teacher' | 'individual';
  study_year?: string;
  department_id?: number;
  street_number: string;
  street_name: string;
  postal_code: string;
  city: string;
}

export interface CreateUserResponse {
  login: string;
  password: string;
  message: string;
}

@Injectable({ providedIn: 'root' })
// API access to create a user and load the departments.
export class CreateUserService {
  private readonly apiUrl = 'http://localhost:8080/api';

  constructor(private http: HttpClient) {}

  getDepartments(): Observable<Department[]> {
    return this.http.get<Department[]>(`${this.apiUrl}/departments`);
  }

  createUser(data: CreateUserPayload): Observable<CreateUserResponse> {
    return this.http.post<CreateUserResponse>(`${this.apiUrl}/users`, data);
  }
}
