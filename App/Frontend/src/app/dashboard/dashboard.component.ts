/**********************************/
/* @since 01/01/2025              */
/* @author K10s Open Source Team  */
/**********************************/

import { AfterViewInit, Component, ElementRef, ViewChild } from '@angular/core';
import { NavComponent } from '../nav/nav.component';
import { FooterComponent } from "../footer/footer.component";
import { TranslatePipe } from "@ngx-translate/core";
import {StatsService} from "../services/stats.service";
import {ByteFormatPipe} from "../byte-format.pipe";
import {Color, NgxChartsModule, ScaleType} from "@swimlane/ngx-charts";

@Component({
    selector: 'app-dashboard',
    templateUrl: './dashboard.component.html',
    styleUrls: ['./dashboard.component.css'],
    imports: [NavComponent, FooterComponent, TranslatePipe, ByteFormatPipe,NgxChartsModule],
    standalone: true
})

export class DashboardComponent implements AfterViewInit {
    onRightClick(event: MouseEvent) {
        event.preventDefault();
    }

    // Fullscreen button
    @ViewChild('dashboardMain') dashboardMain!: ElementRef;
    @ViewChild('dashboardTitle') dashboardTitle!: ElementRef;


    ngAfterViewInit(): void {
        const fullscreenButton = document.getElementById('dashboard-fullscreen-button');
        if (fullscreenButton) fullscreenButton.addEventListener('click', () => this.toggleFullscreen());
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


    // get stats
    usage: any = null;
    memoryChartData: any[] = [];
    cpuChartData: any[] = [];
    diskUsagePercentage: number = 0.00;

    colorScheme:Color = {
        name: 'customScheme',
        selectable: true,
        group: ScaleType.Ordinal,
        domain: []
    };
    colorSchemeCpu: Color = {...this.colorScheme}
    diskColor = "";

    constructor(private usageService: StatsService) {}

    ngOnInit(): void {
        this.usageService.login().subscribe({
            next: () => {
                this.loadUsage();
            },
            error: (error) => {
                console.error(error);
            }
        });
    }


    loadUsage(): void {
        this.usageService.getStats().subscribe({
            next: (data) => {
                console.log(data)
                this.usage = data;
                this.updateChartData();
            },
            error: (error) => {
                console.error(error);
            }
        });
    }

    updateChartData(): void {
        this.memoryChartData = [
            { name: 'Used', value: this.usage.MemUsage },
        ];
        this.cpuChartData = [
            { name: 'Used', value: this.usage?.CpuUsage || 0 },
        ];
        this.diskUsagePercentage = (this.usage.DiskUsage / this.usage.DiskCapacity) * 100;

        this.colorScheme = {
            ...this.colorScheme,
            domain: [this.getUsageColor(this.usage.MemUsage), '#E0E0E0']
        };
        this.colorSchemeCpu = {
            ...this.colorScheme,
            domain: [this.getUsageColor(this.usage.CpuUsage), '#E0E0E0']
        };
        this.diskColor = this.getUsageColor(this.diskUsagePercentage);
    }

    getUsageColor(usage: number): string {
        const rootStyles = getComputedStyle(document.documentElement);
        const green = rootStyles.getPropertyValue('--status-green').trim();
        const orange = rootStyles.getPropertyValue('--status-orange').trim();
        const red = rootStyles.getPropertyValue('--status-red').trim();

        if (usage < 50) return green;
        if (usage < 80) return orange;
        return red;
    }

}