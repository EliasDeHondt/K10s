/**********************************/
/* @since 01/01/2025              */
/* @author K10s Open Source Team  */
/**********************************/

import { Component, AfterViewInit } from '@angular/core';
import { Router, RouterModule } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { TranslatePipe, TranslateService } from "@ngx-translate/core";
import anime from 'animejs/lib/anime.es.js';
import { AuthService } from "../services/auth.service";
import { take, tap } from "rxjs";
import {NotificationService} from "../services/notification.service";

@Component({
    selector: 'app-login',
    standalone: true,
    templateUrl: './login.component.html',
    styleUrls: ['./login.component.css'],
    imports: [FormsModule, RouterModule, TranslatePipe],
})

export class LoginComponent implements AfterViewInit {
    username: string = '';
    password: string = '';

    constructor(private router: Router, private notificationService: NotificationService, private translate: TranslateService, private authService: AuthService) {
        this.translate.setDefaultLang('en');
        this.authService.isLoggedIn().pipe(take(1),
            tap(isAuthenticated => {
                localStorage.setItem('username', this.username);
                if (isAuthenticated) {
                    this.router.navigate(['/dashboard']);
                }
            }))
            .subscribe()
    }

    ngAfterViewInit() {
        const cubes = document.querySelectorAll('g');
        cubes.forEach((cube, index) => {
            const transform = cube.getAttribute('transform') || 'translate(0,0) scale(1)';
            const translateMatch = transform.match(/translate\(([^,]+),([^,]+)\)/);
            const scaleMatch = transform.match(/scale\(([^)]+)\)/);

            const currentTranslateX = translateMatch ? parseFloat(translateMatch[1]) : 0;
            const currentTranslateY = translateMatch ? parseFloat(translateMatch[2]) : 0;
            const scale = scaleMatch ? parseFloat(scaleMatch[1]) : 1;

            anime({
                targets: cube,
                translateY: [currentTranslateY, currentTranslateY - 150],
                translateX: [currentTranslateX, currentTranslateX],
                scale: [scale, scale],
                duration: 1500,
                direction: 'alternate',
                loop: true,
                delay: index * 100,
                endDelay: (el, i, l) => (l - i) * 100
            });
        });
    }

    onSubmit() {
        if (this.username && this.password) {
            this.authService.login(this.username, this.password).subscribe({
                next: () => {
                    this.router.navigate(['/dashboard']);
                }
            })
        } else this.notificationService.showNotification(this.translate.instant('NOTIF.AUTH.INVALID'), 'error');
    }
}