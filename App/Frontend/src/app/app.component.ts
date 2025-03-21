/**********************************/
/* @since 01/01/2025              */
/* @author K10s Open Source Team  */
/**********************************/

import { Component, inject } from '@angular/core';
import { Router, NavigationEnd, RouterModule } from '@angular/router';
import { Title } from '@angular/platform-browser';
import { filter, map } from 'rxjs/operators';
import { CommonModule } from '@angular/common';
import { NotificationComponent } from "./notification/notification.component";

@Component({
    selector: 'app-root',
    standalone: true,
    imports: [RouterModule, CommonModule, NotificationComponent],
    templateUrl: './app.component.html',
    providers: [Title]
})

export class AppComponent {
    private router = inject(Router);
    private titleService = inject(Title);

    constructor() {
        if (this.titleService) {
            this.router.events.pipe(
                filter(event => event instanceof NavigationEnd),
                map(() => {
                    let route = this.router.routerState.root;
                    while (route.firstChild) {
                        route = route.firstChild;
                    }
                    return route.snapshot.data['title'] || 'K10s';
                })
            ).subscribe(title => {
                this.titleService.setTitle(title);
            });
        }
    }
}