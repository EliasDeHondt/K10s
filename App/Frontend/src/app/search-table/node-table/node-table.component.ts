/**********************************/
/* @since 01/01/2025              */
/* @author K10s Open Source Team  */
/**********************************/

import {Component, Input} from '@angular/core';
import {Node} from "../../domain/Kubernetes";
import {TranslatePipe} from "@ngx-translate/core";
import {ScrollService} from "../../services/scroll.service";
import {SkeletonLoaderDirective} from "../skeleton-loader.directive";

@Component({
    selector: 'app-node-table',
    imports: [
        TranslatePipe,
        SkeletonLoaderDirective
    ],
    templateUrl: './node-table.component.html',
    standalone: true,
    styleUrls: ['../../search/search.component.css']
})

export class NodeTableComponent {
    @Input({required: true}) nodes!: Node[];
    @Input({required: true}) isLoading!: boolean;

    constructor(private scrollService: ScrollService) {
    }

    onScroll(): void {
        this.scrollService.emitScroll()
    }

}