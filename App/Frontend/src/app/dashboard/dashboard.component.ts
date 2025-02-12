/**********************************/
/* @since 01/01/2025              */
/* @author K10s Open Source Team  */
/**********************************/

import { AfterViewInit, Component, ElementRef, ViewChild } from '@angular/core';
import { NavComponent } from '../nav/nav.component';
import { FooterComponent } from "../footer/footer.component";
import {TranslatePipe} from "@ngx-translate/core";
import {SpiderWebComponent} from "../spider-web/spider-web.component";

@Component({
    selector: 'app-dashboard',
    templateUrl: './dashboard.component.html',
    styleUrls: ['./dashboard.component.css'],
    imports: [NavComponent, FooterComponent, TranslatePipe, SpiderWebComponent],
    standalone: true
})

export class DashboardComponent implements AfterViewInit {
    // Fullscreen button
    @ViewChild('dashboardMain') dashboardMain!: ElementRef;
    @ViewChild('dashboardTitle') dashboardTitle!: ElementRef;

    constructor() {}

    ngAfterViewInit(): void {
        const fullscreenButton = document.getElementById('dashboard-fullscreen-button');
        if (fullscreenButton) {
        fullscreenButton.addEventListener('click', () => this.toggleFullscreen());
        }
    }

    toggleFullscreen(): void {
        const dashboardHtml = document.documentElement;
        const dashboardMainEl = this.dashboardMain.nativeElement;
        const dashboardTitleEl = this.dashboardTitle.nativeElement;

        if (!document.fullscreenElement) {
        if (dashboardHtml.requestFullscreen) dashboardHtml.requestFullscreen();
        else if ((dashboardHtml as any).mozRequestFullScreen) (dashboardHtml as any).mozRequestFullScreen();
        else if ((dashboardHtml as any).webkitRequestFullscreen) (dashboardHtml as any).webkitRequestFullscreen();
        else if ((dashboardHtml as any).msRequestFullscreen) (dashboardHtml as any).msRequestFullscreen();

        dashboardMainEl.classList.add('fullscreen');
        dashboardMainEl.style.gridArea = '1 / 1 / -1 / -1';
        dashboardTitleEl.style.display = 'block';
        } else {
        if (document.exitFullscreen) document.exitFullscreen();
        else if ((document as any).mozCancelFullScreen) (document as any).mozCancelFullScreen();
        else if ((document as any).webkitExitFullscreen) (document as any).webkitExitFullscreen();
        else if ((document as any).msExitFullscreen) (document as any).msExitFullscreen();

        dashboardMainEl.classList.remove('fullscreen');
        dashboardMainEl.style.gridArea = '';
        dashboardTitleEl.style.display = 'none';
        }
    }
}