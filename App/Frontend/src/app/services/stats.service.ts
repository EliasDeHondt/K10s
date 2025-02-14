import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class StatsService {
  private apiUrl = 'http://localhost:8080/secured/stats';

  constructor(private http: HttpClient) {}

  getStats(): Observable<any> {
    return this.http.get<any>(this.apiUrl);
  }
}
