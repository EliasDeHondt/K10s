import {effect, Injectable} from "@angular/core";
import {environment} from "../../environments/environment";
import {HttpClient} from "@angular/common/http";
import {Observable} from "rxjs";
import {VisualizationData} from "../spider-web/spider-web.component";

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

    getVisualization(): Observable<VisualizationData> {
        return this.http.get<VisualizationData>(this.apiUrl, {withCredentials: true});
    }
}