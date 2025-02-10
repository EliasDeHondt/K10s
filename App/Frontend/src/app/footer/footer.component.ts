/**********************************/
/* @since 01/01/2025              */
/* @author K10s Open Source Team  */
/**********************************/

import { Component, OnInit } from '@angular/core';
import {TranslatePipe} from "@ngx-translate/core";

@Component({
    selector: 'app-footer',
    templateUrl: './footer.component.html',
    standalone: true,
    imports: [TranslatePipe],
    styleUrls: ['./footer.component.css']
})

export class FooterComponent implements OnInit {
    currentYear: number = new Date().getFullYear();

    ngOnInit(): void {
        this.currentYear = new Date().getFullYear();
    }
}