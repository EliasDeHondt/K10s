/**********************************/
/* @since 01/01/2025              */
/* @author K10s Open Source Team  */
/**********************************/

import { Component, Input } from '@angular/core';
import { TranslatePipe } from "@ngx-translate/core";
import { Secret } from "../../domain/Kubernetes";

@Component({
    selector: 'app-secret-table',
    imports: [
        TranslatePipe,
    ],
    templateUrl: './secret-table.component.html',
    standalone: true,
    styleUrls: ['../../search/search.component.css', './secret-table.component.css']
})

export class SecretTableComponent {
    @Input({required: true}) secrets!: Secret[];

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