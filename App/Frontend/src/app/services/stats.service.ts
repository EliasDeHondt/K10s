/**********************************/
/* @since 01/01/2025              */
/* @author K10s Open Source Team  */
/**********************************/

import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { Metrics } from "../domain/Metrics";
import { environment } from "../../environments/environment";

@Injectable({
    providedIn: 'root'
})

export class StatsService {
    private apiUrl = `${environment.BASE_URL}/secured/stats`;

    constructor(private http: HttpClient) { }

    getStats(): Observable<Metrics> {
        return this.http.get<Metrics>(this.apiUrl, {withCredentials: true});
    }
}