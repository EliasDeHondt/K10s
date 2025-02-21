/**********************************/
/* @since 01/01/2025              */
/* @author K10s Open Source Team  */
/**********************************/

import { Component, Input } from '@angular/core';
import { Node } from "../../domain/Kubernetes";
import { TranslatePipe } from "@ngx-translate/core";

@Component({
    selector: 'app-node-table',
    imports: [
        TranslatePipe
    ],
    templateUrl: './node-table.component.html',
    standalone: true,
    styleUrls: ['../../search/search.component.css']
})

export class NodeTableComponent {
    @Input({required: true}) nodes!: Node[];
}