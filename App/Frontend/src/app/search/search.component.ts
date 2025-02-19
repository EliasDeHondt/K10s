/**********************************/
/* @since 01/01/2025              */
/* @author K10s Open Source Team  */
/**********************************/

import { Component } from '@angular/core';
import { NavComponent } from '../nav/nav.component';
import { FooterComponent } from "../footer/footer.component";
import { CommonModule } from "@angular/common";
import { TranslatePipe } from "@ngx-translate/core";
import { LoadingComponent } from "../loading/loading.component";

@Component({
    selector: 'app-search',
    templateUrl: './search.component.html',
    styleUrls: ['./search.component.css'],
    imports: [NavComponent, FooterComponent, CommonModule, TranslatePipe],
    standalone: true
})

export class SearchComponent {
    isLoading: boolean = true;
    dropdowns: { [key: string]: boolean } = {
        searchDropdown1: false,
        searchDropdown2: false,
        searchDropdown3: false
    };
    selectedNode: string = 'None';
    selectedNamespace: string = 'None';
    searchResults: any[] = [];

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

    ngOnInit(): void {
        this.getData();
    }

    getData(): void {
        this.isLoading = true;
        //todo get
    }
}