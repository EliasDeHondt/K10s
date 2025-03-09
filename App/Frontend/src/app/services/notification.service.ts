/**********************************/
/* @since 01/01/2025              */
/* @author K10s Open Source Team  */
/**********************************/

import { Injectable } from '@angular/core';
import { Subject } from 'rxjs';

@Injectable({
    providedIn: 'root'
})

export class NotificationService {
    private notificationSubject = new Subject<{ message: string; type: 'success' | 'info' | 'error' }>();
    notification$ = this.notificationSubject.asObservable();

    showNotification(message: string, type: 'success' | 'info' | 'error' = 'success') {
        this.notificationSubject.next({ message, type });
    }
}