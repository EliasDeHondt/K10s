/**********************************/
/* @since 01/01/2025              */
/* @author K10s Open Source Team  */
/**********************************/

import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import {catchError, Observable, throwError} from 'rxjs';
import {BASE_URL} from "./stats.service";

@Injectable({
    providedIn: 'root'
})

export class DeploymentService {
    private uploadUrl = `${BASE_URL}/upload`;

    constructor(private http: HttpClient) {}

    uploadYaml(yaml: string): Observable<any> {
        return this.http.post<any>(this.uploadUrl, { yaml }, { withCredentials: true }).pipe(
            catchError(error => {
                console.error('YAML upload failed', error);
                return throwError(() => error);
            })
        );
    }
}