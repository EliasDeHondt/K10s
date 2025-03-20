/**********************************/
/* @since 01/01/2025              */
/* @author K10s Open Source Team  */
/**********************************/

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
        document.addEventListener('visibilitychange', this.handleVisibilityChange);
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
            }
        }

        this.socket.onerror = () => {
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
        document.removeEventListener('visibilitychange', this.handleVisibilityChange);
    }

    isConnected() {
        return this.socket && this.socket.readyState === WebSocket.OPEN;
    }

    private handleVisibilityChange = () => {
        if (document.visibilityState === 'visible' && !this.isConnected()) {
            this.connect();
        }
    };
}