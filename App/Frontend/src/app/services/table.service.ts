/**********************************/
/* @since 01/01/2025              */
/* @author K10s Open Source Team  */
/**********************************/

import {effect, Injectable, signal} from "@angular/core";
import {HttpClient} from "@angular/common/http";
import {environment} from "../../environments/environment";
import {PaginatedResponse} from "../domain/Kubernetes";

@Injectable({
    providedIn: 'root'
})

export class TableService {
    private tableUrl = `${environment.BASE_URL}/secured/table`;
    element = signal('')
    namespace = signal('')
    data = signal<PaginatedResponse>({Response: [], PageToken: ''})

    constructor(private http: HttpClient) {
        effect(() => {
            this.getTable(this.element(), this.namespace())
        });
    }

    setElement(element: string) {
        this.element.set(element);
    }

    setNamespace(name: string) {
        this.namespace.set(name);
    }

    private getTable(element: string, namespace: string) {
        if (!element) {
            return
        }

        this.http.get<PaginatedResponse>(this.tableUrl + `?element=${element}&namespace=${namespace}`, {withCredentials: true}).subscribe(data => {
            this.data.set(data)
        })
    }
}