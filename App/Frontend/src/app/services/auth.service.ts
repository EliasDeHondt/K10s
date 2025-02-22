/**********************************/
/* @since 01/01/2025              */
/* @author K10s Open Source Team  */
/**********************************/

import {Injectable} from '@angular/core';
import {HttpClient, HttpErrorResponse} from '@angular/common/http';
import {catchError, Observable, of, throwError} from 'rxjs';
import {environment} from "../../environments/environment";
import {map} from "rxjs/operators";

@Injectable({
    providedIn: 'root'
})

export class AuthService {
    private loginUrl = `${environment.BASE_URL}/login`;
    private logoutUrl = `${environment.BASE_URL}/logout`;
    private isloggedInUrl = `${environment.BASE_URL}/isloggedin`;

    constructor(private http: HttpClient) {
    }

    isLoggedIn(): Observable<boolean> {
        return this.http.get<boolean>(this.isloggedInUrl, {withCredentials: true, observe: 'response'}).pipe(
            map(response => {
                return response.body!;
            }),
            catchError((error: HttpErrorResponse) => {
                return of(false);
            })
        )
    }

    login(username: string, password: string): Observable<string> {
        const loginData = {username: username, password: password};

        return this.http.post<string>(this.loginUrl, loginData, {withCredentials: true}).pipe(
            catchError(error => {
                const statusCode = error.status;
                if (statusCode === 401) {
                    alert("Please enter valid credentials")
                }

                return throwError(() => error);
            })
        );
    }

    logout(): Observable<string> {
        return this.http.get<string>(this.logoutUrl, {withCredentials: true}).pipe(
            catchError(error => {
                return throwError(() => error);
            })
        )
    }
}