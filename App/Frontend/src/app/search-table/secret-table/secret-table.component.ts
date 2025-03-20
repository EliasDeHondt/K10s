/**********************************/
/* @since 01/01/2025              */
/* @author K10s Open Source Team  */
/**********************************/

import {Component, Input} from '@angular/core';
import {TranslatePipe} from "@ngx-translate/core";
import {Secret} from "../../domain/Kubernetes";
import {SearchTooltipService} from "../../services/tooltip.service";
import {ScrollService} from "../../services/scroll.service";

@Component({
    selector: 'app-secret-table',
    imports: [
        TranslatePipe,
    ],
    templateUrl: './secret-table.component.html',
    standalone: true,
    styleUrls: ['../../search/search.component.css', '../config-map-table/config-map-table.component.css'],
})

export class SecretTableComponent {
    @Input({required: true}) secrets!: Secret[];

    constructor(private tooltipService: SearchTooltipService, private scrollService: ScrollService) {
    }

    onScroll(): void {
        this.scrollService.emitScroll()
    }

    protected readonly Object = Object;

    showTooltip(event: MouseEvent, data: Record<string, any> | undefined) {
        this.tooltipService.showTooltip(event, data);
    }

    hideTooltip() {
        this.tooltipService.hideTooltip();
    }
}