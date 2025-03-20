/**********************************/
/* @since 01/01/2025              */
/* @author K10s Open Source Team  */
/**********************************/

import {Component, inject, OnInit} from '@angular/core';
import {NavComponent} from '../nav/nav.component';
import {FooterComponent} from "../footer/footer.component";
import {CommonModule} from "@angular/common";
import {TranslatePipe} from "@ngx-translate/core";
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
import {SecretTableComponent} from "../search-table/secret-table/secret-table.component";
import {SecretCastPipe} from "../pipes/secret-cast.pipe";
import {FilterDataService} from "../services/filterdata.service";
import {Namespace} from "../domain/Kubernetes";
import {FormsModule} from "@angular/forms";
import {LoadingService} from "../services/loading.service";
import {ScrollService} from "../services/scroll.service";

@Component({
    selector: 'app-search',
    templateUrl: './search.component.html',
    styleUrls: ['./search.component.css'],
    imports: [NavComponent, FooterComponent, CommonModule, TranslatePipe, PodCastPipe, PodTableComponent, NodeTableComponent, NodeCastPipe, ServiceTableComponent, ServiceCastPipe, DeploymentTableComponent, DeploymentCastPipe, ConfigMapTableComponent, ConfigMapCastPipe, SecretTableComponent, SecretCastPipe, FormsModule],
    standalone: true
})

export class SearchComponent implements OnInit {
    tableService = inject(TableService);
    filterDataService = inject(FilterDataService);
    loadingService = inject(LoadingService);
    dropdowns: { [key: string]: boolean } = {
        searchDropdown1: false,
        searchDropdown2: false,
        searchDropdown3: false
    };
    namespaces: Namespace[] = [];
    nodeNames: string[] = [];

    pageSize: number = 20;

    ngOnInit(): void {
        this.getNamespaces()
        this.getNodeNames()
        this.tableService.getTable(this.tableService.element(), this.tableService.namespace(), this.tableService.node())
        this.scrollService.scroll$.subscribe(() => {
            this.onScroll()
        })
    }

    constructor(private scrollService: ScrollService) {
        this.tableService.setElement('nodes');
    }

    setPageSize(size: number) {
        this.pageSize = size;
        this.loadingService.isLoading.set(true);
        this.tableService.getTable(
            this.tableService.element(),
            this.tableService.namespace(),
            this.tableService.node(),
            this.pageSize
        );
    }

    updateElement(filter: string) {
        this.loadingService.isLoading.set(true);
        this.tableService.setElement(filter);
        this.tableService.getTable(filter, this.tableService.namespace(), this.tableService.node(), this.pageSize);
    }

    selectNode(node: string) {
        this.loadingService.isLoading.set(true);
        this.tableService.setNodeName(node);
        this.tableService.getTable(this.tableService.element(), this.tableService.namespace(), node, this.pageSize);
        this.toggleDropdown('searchDropdown1');
    }

    selectNamespace(namespace: string) {
        this.loadingService.isLoading.set(true);
        this.tableService.setNamespace(namespace);
        this.tableService.getTable(this.tableService.element(), namespace, this.tableService.node(), this.pageSize);
        this.toggleDropdown('searchDropdown2');
    }

    toggleDropdown(dropdownKey: string) {
        for (let key in this.dropdowns) {
            this.dropdowns[key] = key === dropdownKey ? !this.dropdowns[key] : false;
        }
    }


    getNamespaces() {
        this.filterDataService.getNamespaces().subscribe(response => {
            this.namespaces = response;
        })
    }

    getNodeNames() {
        this.filterDataService.getNodeNames().subscribe(response => {
            this.nodeNames = response;
        })
    }

    isDefaultNodeSelected(): boolean {
        return this.tableService.element() === 'nodes';
    }

    onScroll() {
        const container = document.querySelector('.search-table');
        if (container) {
            const {scrollTop, scrollHeight, clientHeight} = container;
            const scrollPosition = scrollTop + clientHeight;

            if (scrollPosition >= scrollHeight / 2 && !this.loadingService.isLoading()) {
                this.tableService.getNextPage(this.tableService.element(), this.tableService.namespace(), this.tableService.node(), this.pageSize);
            }
        }
    }

}