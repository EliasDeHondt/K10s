import {effect, Injectable} from "@angular/core";
import {environment} from "../../environments/environment";
import {HttpClient} from "@angular/common/http";
import {Observable} from "rxjs";
import {Visualization} from "../domain/Visualization";

@Injectable({
    providedIn: 'root'
})

export class VisualizationService {
    private apiUrl = `${environment.BASE_URL}/secured/visualization`;

    constructor(private http: HttpClient) {
        effect(() => {
            this.getVisualization();
        });
    }

    getVisualization(): Observable<Visualization> {
        return this.http.get<Visualization>(this.apiUrl, {withCredentials: true});
    }
}