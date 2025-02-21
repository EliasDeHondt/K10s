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
    data = signal<PaginatedResponse>({Response: [], PageToken: ''})

    constructor(private http: HttpClient) {
        effect(() => {
            this.getTable(this.element())
        });
    }

    setElement(element: string){
        this.element.set(element);
    }

    private getTable(element: string){
        if(!element){
            return
        }

        this.http.get<PaginatedResponse>(this.tableUrl + `?element=${element}`, {withCredentials: true}).subscribe(data => {
            this.data.set(data)
        })
    }

}