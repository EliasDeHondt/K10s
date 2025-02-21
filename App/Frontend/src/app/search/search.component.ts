/**********************************/
/* @since 01/01/2025              */
/* @author K10s Open Source Team  */
/**********************************/

import {Component, inject, signal} from '@angular/core';
import { NavComponent } from '../nav/nav.component';
import { FooterComponent } from "../footer/footer.component";
import { CommonModule } from "@angular/common";
import { TranslatePipe } from "@ngx-translate/core";
import {PodTableComponent} from "../search-table/pod-table/pod-table.component";
import {TableService} from "../services/table.service";
import {PodCastPipe} from "../pipes/pod-cast.pipe";
import {NodeTableComponent} from "../search-table/node-table/node-table.component";
import {NodeCastPipe} from "../pipes/node-cast.pipe";
import {ServiceTableComponent} from "../search-table/service-table/service-table.component";
import {ServiceCastPipe} from "../pipes/service-cast.pipe";
import {DeploymentTableComponent} from "../search-table/deployment-table/deployment-table.component";
import {DeploymentCastPipe} from "../pipes/deployment-cast.pipe";
import {ConfigMapTableComponent} from "../search-table/config-map-table/config-map-table.component";
import {ConfigMapCastPipe} from "../pipes/config-map-cast.pipe";

@Component({
    selector: 'app-search',
    templateUrl: './search.component.html',
    styleUrls: ['./search.component.css'],
    imports: [NavComponent, FooterComponent, CommonModule, TranslatePipe, PodCastPipe, PodTableComponent, NodeTableComponent, NodeCastPipe, ServiceTableComponent, ServiceCastPipe, DeploymentTableComponent, DeploymentCastPipe, ConfigMapTableComponent, ConfigMapCastPipe],
    standalone: true
})

export class SearchComponent {
    tableService = inject(TableService);
    isLoading = signal(false)
    dropdowns: { [key: string]: boolean } = {
        searchDropdown1: false,
        searchDropdown2: false,
        searchDropdown3: false
    };
    selectedNode: string = 'None';
    selectedNamespace: string = 'None';

    updateElement(filter: string) {
        this.isLoading.set(true);
        this.tableService.setElement(filter)
        this.isLoading.set(false);
    }

    toggleDropdown(dropdownKey: string) {
        for (let key in this.dropdowns) {
            this.dropdowns[key] = key === dropdownKey ? !this.dropdowns[key] : false;
        }
    }

    selectNode(node: string) {
        this.selectedNode = node;
        this.toggleDropdown('searchDropdown1');
    }

    selectNamespace(namespace: string) {
        this.selectedNamespace = namespace;
        this.toggleDropdown('searchDropdown2');
    }

}