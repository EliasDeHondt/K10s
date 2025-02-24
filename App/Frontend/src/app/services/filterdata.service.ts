/**********************************/
/* @since 01/01/2025              */
/* @author K10s Open Source Team  */
/**********************************/

import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';
import {environment} from "../../environments/environment";
import {Namespace} from "../domain/Kubernetes";

@Injectable({
    providedIn: 'root'
})

export class FilterDataService {
    private namespacesUrl = `${environment.BASE_URL}/secured/namespaces`;

    constructor(private http: HttpClient) {
    }

    getNamespaces(): Observable<Namespace[]> {
        return this.http.get<Namespace[]>(this.namespacesUrl, {withCredentials: true});
    }
}