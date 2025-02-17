import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import {catchError, Observable, tap, throwError} from 'rxjs';
export const BASE_URL = 'http://localhost:8080';

@Injectable({
  providedIn: 'root'
})
export class StatsService {
  private apiUrl = `${BASE_URL}/secured/stats`;
  private loginUrl = `${BASE_URL}/login`;

  constructor(private http: HttpClient) {}

  login(): Observable<any> {
    const loginData = { username: 'admin', password: 'password' };

    return this.http.post<any>(this.loginUrl, loginData,{ withCredentials: true }).pipe(
        catchError(error => {
          return throwError(() => error);
        })
    );
  }

  getStats(): Observable<any> {
    return this.http.get<any>(this.apiUrl,{withCredentials: true});
  }
}
