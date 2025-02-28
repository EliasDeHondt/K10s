/**********************************/
/* @since 01/01/2025              */
/* @author K10s Open Source Team  */
/**********************************/

import { Component, Input } from '@angular/core';
import { TranslatePipe } from "@ngx-translate/core";
import { ConfigMap } from "../../domain/Kubernetes";
import { SearchTooltipService } from "../../services/tooltip.service";

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

    constructor(private tooltipService: SearchTooltipService) {}

    protected readonly Object = Object;

    showTooltip(event: MouseEvent, data: Record<string, any> | undefined) {
        this.tooltipService.showTooltip(event, data);
    }

    hideTooltip() {
        this.tooltipService.hideTooltip();
    }
}
