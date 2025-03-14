/**********************************/
/* @since 01/01/2025              */
/* @author K10s Open Source Team  */
/**********************************/

import { AfterViewInit, Component, ElementRef, OnDestroy, OnInit, ViewChild } from '@angular/core';
import { NavComponent } from '../nav/nav.component';
import { FooterComponent } from '../footer/footer.component';
import { TranslatePipe } from '@ngx-translate/core';
import { ByteFormatPipe } from '../byte-format.pipe';
import { Color, NgxChartsModule, ScaleType } from '@swimlane/ngx-charts';
import { SpiderWebComponent } from '../spider-web/spider-web.component';
import { StatWebSocketService } from '../services/statWebsocket.service';
import { Metrics } from '../domain/Metrics';

@Component({
    selector: 'app-dashboard',
    templateUrl: './dashboard.component.html',
    styleUrls: ['./dashboard.component.css'],
    imports: [NavComponent, FooterComponent, TranslatePipe, ByteFormatPipe, NgxChartsModule, SpiderWebComponent],
    standalone: true,
})
export class DashboardComponent implements AfterViewInit, OnInit, OnDestroy {
    @ViewChild('dashboardMain') dashboardMain!: ElementRef;
    @ViewChild('dashboardTitle') dashboardTitle!: ElementRef;
    @ViewChild('spiderWeb', { static: false }) spiderWeb!: ElementRef;

    usage: Metrics | undefined = undefined;
    memoryChartData: any[] = [];
    cpuChartData: any[] = [];
    diskUsagePercentage: number = 0.0;
    diskUsage: number = 0.0;
    diskCapacity: number = 0.0;


    colorScheme: Color = {
        name: 'customScheme',
        selectable: true,
        group: ScaleType.Ordinal,
        domain: [],
    };
    colorSchemeCpu: Color = { ...this.colorScheme };
    diskColor = '';
    loading: boolean = false;

    constructor(private usageService: StatWebSocketService) {}

    ngOnInit(): void {
        this.usageService.connect();

        this.usageService.getMetrics().subscribe({
            next: (data) => {
                this.updateChartData(data);
                this.loading = false;
            },
            error: (error) => {
                console.error(error);
                this.loading = true;
            },
        });
    }

    ngAfterViewInit(): void {
        const fullscreenButton = document.getElementById('dashboard-fullscreen-button');
        if (fullscreenButton) {
            fullscreenButton.addEventListener('click', () => this.toggleFullscreen());
        }

        const downloadButton = document.getElementById('dashboard-download-button');
        if (downloadButton) {
            downloadButton.addEventListener('click', () => this.downloadSpiderWebAsPNG());
        }
    }

    downloadSpiderWebAsPNG(): void {
        const spiderWebSvg = document.querySelector('app-spider-web svg');
        if (!spiderWebSvg) return;

        const svgClone = spiderWebSvg.cloneNode(true) as SVGSVGElement;
        const imageElements = svgClone.querySelectorAll('image[href]') as NodeListOf<SVGImageElement>;
        const promises: Promise<void>[] = [];

        imageElements.forEach((img: SVGImageElement) => {
            const href = img.getAttribute('href') || img.getAttribute('xlink:href');
            if (href) {
                promises.push(
                    fetch(href)
                        .then((response) => response.blob())
                        .then((blob) => {
                            return new Promise<void>((resolve) => {
                                const reader = new FileReader();
                                reader.onloadend = () => {
                                    img.setAttribute('href', reader.result as string);
                                    resolve();
                                };
                                reader.readAsDataURL(blob);
                            });
                        })
                        .catch((error) => {
                            console.error(`Failed to load image at ${href}:`, error);
                        })
                );
            }
        });

        Promise.all(promises).then(() => {
            const serializer = new XMLSerializer();
            const svgString = serializer.serializeToString(svgClone);
            const svgBlob = new Blob([svgString], { type: 'image/svg+xml' });
            const svgUrl = URL.createObjectURL(svgBlob);

            const img = new Image();
            img.onload = () => {
                const canvas = document.createElement('canvas');
                canvas.width = img.width;
                canvas.height = img.height;
                const ctx = canvas.getContext('2d');
                if (ctx) {
                    ctx.drawImage(img, 0, 0);
                    const pngUrl = canvas.toDataURL('image/png');
                    const a = document.createElement('a');
                    a.href = pngUrl;
                    a.download = 'spider-web.png';
                    document.body.appendChild(a);
                    a.click();
                    document.body.removeChild(a);
                }
                URL.revokeObjectURL(svgUrl);
            };
            img.src = svgUrl;
        });
    }

    ngOnDestroy(): void {
        this.usageService.disconnect();
        const fullscreenButton = document.getElementById('dashboard-fullscreen-button');
        if (fullscreenButton) {
            fullscreenButton.removeEventListener('click', () => this.toggleFullscreen());
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

    valueFormatting(usage: number): string {
        return usage.toFixed(2) + '%';
    }


    updateChartData(metrics: Metrics): void {
        this.memoryChartData = [{ name: 'Used', value: parseFloat(metrics.MemUsage.toFixed(2)) }];
        this.cpuChartData = [{ name: 'Used', value: parseFloat(metrics.CpuUsage.toFixed(2)) }];
        this.diskUsage = metrics.DiskUsage;
        this.diskCapacity = metrics.DiskCapacity;

        this.diskUsagePercentage = (metrics.DiskUsage / metrics.DiskCapacity) * 100;

        this.colorScheme = {
            ...this.colorScheme,
            domain: [this.getUsageColor(metrics.MemUsage), '#E0E0E0'],
        };
        this.colorSchemeCpu = {
            ...this.colorScheme,
            domain: [this.getUsageColor(metrics.CpuUsage), '#E0E0E0'],
        };
        this.diskColor = this.getUsageColor(this.diskUsagePercentage);
    }

    getUsageColor(usage: number): string {
        const rootStyles = getComputedStyle(document.documentElement);
        const green = rootStyles.getPropertyValue('--status-green').trim();
        const orange = rootStyles.getPropertyValue('--status-orange').trim();
        const red = rootStyles.getPropertyValue('--status-red').trim();

        if (usage < 55) return green;
        if (usage < 85) return orange;
        return red;
    }
}
