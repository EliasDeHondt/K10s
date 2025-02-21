/**********************************/
/* @since 01/01/2025              */
/* @author K10s Open Source Team  */
/**********************************/

import { Component, Input } from '@angular/core';
import { TranslatePipe } from "@ngx-translate/core";
import { ConfigMap } from "../../domain/Kubernetes";

@Component({
    selector: 'app-config-map-table',
    imports: [
        TranslatePipe
    ],
    templateUrl: './config-map-table.component.html',
    standalone: true,
    styleUrls: ['../../search/search.component.css','./config-map-table.component.css']
})

export class ConfigMapTableComponent {
    @Input({ required: true }) configMaps!: ConfigMap[];

    protected readonly Object = Object;

    showTooltip(event: MouseEvent, data: Record<string, any>) {
        const tooltip = document.getElementById('tooltip');
        if (tooltip) {
        tooltip.textContent = Object.entries(data)
            .map(([key, value]) => `${key}: ${value}`)
            .join('\n');

        tooltip.style.display = 'block';
        tooltip.style.left = `${event.pageX + 10}px`;
        tooltip.style.top = `${event.pageY + 10}px`;
        }
    }

    hideTooltip() {
        const tooltip = document.getElementById('tooltip');
        if (tooltip) {
        tooltip.style.display = 'none';
        }
    }
}
