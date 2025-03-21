/**********************************/
/* @since 01/01/2025              */
/* @author K10s Open Source Team  */
/**********************************/

import {Component, Input} from '@angular/core';
import {TranslatePipe} from "@ngx-translate/core";
import {Deployment} from "../../domain/Kubernetes";
import {ScrollService} from "../../services/scroll.service";
import {SkeletonLoaderDirective} from "../skeleton-loader.directive";

@Component({
    selector: 'app-deployment-table',
    imports: [
        TranslatePipe,
        SkeletonLoaderDirective
    ],
    templateUrl: './deployment-table.component.html',
    standalone: true,
    styleUrls: ['../../search/search.component.css']
})

export class DeploymentTableComponent {
    @Input({required: true}) deployments!: Deployment[];
    @Input({required: true}) isLoading!: boolean;

    constructor(private scrollService: ScrollService) {
    }

    onScroll(): void {
        this.scrollService.emitScroll()
    }

}