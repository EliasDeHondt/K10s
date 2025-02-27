/**********************************/
/* @since 01/01/2025              */
/* @author K10s Open Source Team  */
/**********************************/

import { Component, Input } from '@angular/core';
import { Service } from "../../domain/Kubernetes";
import { TranslatePipe } from "@ngx-translate/core";

@Component({
    selector: 'app-service-table',
    imports: [
        TranslatePipe
    ],
    templateUrl: './service-table.component.html',
    standalone: true,
    styleUrls: ['../../search/search.component.css']
})

export class ServiceTableComponent {
    @Input({ required: true }) services!: Service[];
}