/**********************************/
/* @since 01/01/2025              */
/* @author K10s Open Source Team  */
/**********************************/

import { Component, OnInit, OnDestroy } from '@angular/core';
import { NotificationService } from '../services/notification.service';
import { Subscription } from 'rxjs';
import { NgClass } from '@angular/common';

@Component({
  selector: 'app-notification',
  standalone: true,
  templateUrl: './notification.component.html',
  styleUrls: ['./notification.component.css'],
  imports: [NgClass]
})
export class NotificationComponent implements OnInit, OnDestroy {
  isVisible: boolean = false;
  message: string = '';
  notificationType: string = 'success';
  private subscription!: Subscription;

  constructor(private notificationService: NotificationService) {}

  ngOnInit() {
    this.subscription = this.notificationService.notification$.subscribe({
      next: ({ message, type }) => {
        this.message = message;
        this.notificationType = type;
        this.isVisible = true;
        setTimeout(() => {
          this.isVisible = false;
        }, 2500);
      },
      error: (error) => console.error('Subscription error:', error)
    });
  }

  ngOnDestroy() {
    if (this.subscription) {
      this.subscription.unsubscribe();
    }
  }
}