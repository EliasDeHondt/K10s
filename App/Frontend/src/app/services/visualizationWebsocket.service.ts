/**********************************/
/* @since 01/01/2025              */
/* @author K10s Open Source Team  */
/**********************************/

import {Injectable, OnDestroy} from "@angular/core";
import {Subject} from "rxjs";
import {environment} from "../../environments/environment";
import {Visualization} from "../domain/Visualization";

@Injectable({
    providedIn: 'root'
})

export class VisualizationWebSocketService implements OnDestroy {
    private url = `${environment.BASE_URL}/secured/visualization`;
    private socket!: WebSocket;
    private messagesSubject: Subject<Visualization> = new Subject()

    constructor() {
        document.addEventListener('visibilitychange', this.handleVisibilityChange);
    }

    connect(): void {
        this.socket = new WebSocket(this.url);

        this.socket.onopen = () => {
        }

        this.socket.onmessage = (event: MessageEvent) => {
            try {
                const data: Visualization = JSON.parse(event.data);
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

    getVisualization() {
        return this.messagesSubject.asObservable();
    }

    isConnected() {
        return this.socket && this.socket.readyState === WebSocket.OPEN;
    }

    sendNamespaceFilter(namespace: string) {
        this.socket.send(namespace);
    }

    ngOnDestroy() {
        this.disconnect();
        document.removeEventListener('visibilitychange', this.handleVisibilityChange);
    }

    private handleVisibilityChange = () => {
        if (document.visibilityState === 'visible' && !this.isConnected()) {
            this.connect();
        }
    };

    formatMemory(memory: string | undefined | null): string {
        if (!memory) return '0KB';
        const match = memory.match(/^(\d+)([KMGTP]i?)?$/);
        if (!match) return memory;


        const value = parseInt(match[1], 10);
        const unit = match[2] || '';

        let bytes: number;
        switch (unit.toLowerCase()) {
            case 'ki':
                bytes = value * 1024;
                break;
            case 'mi':
                bytes = value * 1024 * 1024;
                break;
            case 'gi':
                bytes = value * 1024 * 1024 * 1024;
                break;
            case 'ti':
                bytes = value * 1024 * 1024 * 1024 * 1024;
                break;
            case 'pi':
                bytes = value * 1024 * 1024 * 1024 * 1024 * 1024;
                break;
            case '':
                bytes = value;
                break;
            default:
                return memory;
        }

        const gb = bytes / (1000 * 1000 * 1000);
        if (gb >= 1) return `${gb.toFixed(2)} GB`;
        const mb = bytes / (1000 * 1000);
        if (mb >= 1) return `${mb.toFixed(2)} MB`;
        const kb = bytes / 1000;
        return `${kb.toFixed(2)} KB`;
    }
}