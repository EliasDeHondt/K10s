import {Injectable, OnDestroy} from "@angular/core";
import {Subject} from "rxjs";
import {Metrics} from "../domain/Metrics";
import {environment} from "../../environments/environment";

@Injectable({
    providedIn: 'root'
})
export class StatWebSocketService implements OnDestroy {
    private url = `${environment.BASE_URL}/secured/statsocket`;
    private socket!: WebSocket;
    private messagesSubject: Subject<Metrics> = new Subject()

    constructor() {
    }

    connect(): void {
        this.socket = new WebSocket(this.url);

        this.socket.onopen = () => {
        }

        this.socket.onmessage = (event: MessageEvent) => {
            try {
                const data: Metrics = JSON.parse(event.data);
                this.messagesSubject.next(data)
            } catch (error) {
                console.error("[WebSocketService] Error parsing data:", error)
            }
        }

        this.socket.onerror = () => {
            console.error("[WebSocketService] WebSocket Error: " + this.url);
        }

        this.socket.onclose = () => {
        }
    }

    disconnect(): void {
        if (this.socket && this.socket.readyState === WebSocket.OPEN) {
            this.socket.close();
        }
    }

    getMetrics() {
        return this.messagesSubject.asObservable();
    }

    ngOnDestroy() {
        this.disconnect();
    }

}