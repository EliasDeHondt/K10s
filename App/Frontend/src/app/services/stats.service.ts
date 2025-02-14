import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import {catchError, Observable, tap, throwError} from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class StatsService {
  private apiUrl = 'http://localhost:8080/secured/stats';
  private loginUrl = 'http://localhost:8080/login';

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
