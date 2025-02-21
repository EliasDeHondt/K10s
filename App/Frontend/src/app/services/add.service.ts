/**********************************/
/* @since 01/01/2025              */
/* @author K10s Open Source Team  */
/**********************************/

import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { catchError, Observable, throwError } from 'rxjs';
import {environment} from "../../environments/environment";

@Injectable({
    providedIn: 'root'
})

export class AddService {
    private uploadUrl = `${environment.BASE_URL}/secured/createresources`;

    constructor(private http: HttpClient) {}

    uploadYaml(yaml: string): Observable<any> {
        return this.http.post<any>(this.uploadUrl, yaml, {
            headers: { 'Content-Type': 'application/x-yaml' },
            withCredentials: true
        }).pipe(
            catchError(error => {
                console.error('YAML upload failed', error);
                return throwError(() => error);
            })
        );
    }
}