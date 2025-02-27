/**********************************/
/* @since 01/01/2025              */
/* @author K10s Open Source Team  */
/**********************************/

import { Component, Input } from '@angular/core';
import { TranslatePipe } from "@ngx-translate/core";
import { Secret } from "../../domain/Kubernetes";
import {hideTooltip, showTooltip} from "../search-table-util";

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

    showTooltip(event: MouseEvent, data: Record<string, any> | undefined) {
        showTooltip(event, data);
    }

    hideTooltip() {
        hideTooltip();
    }
}