/**********************************/
/* @since 01/01/2025              */
/* @author K10s Open Source Team  */
/**********************************/

import { Component } from '@angular/core';
import { FooterComponent } from "../footer/footer.component";
import { TranslatePipe, TranslateService } from "@ngx-translate/core";

@Component({
    selector: 'app-notfound',
    templateUrl: './notfound.component.html',
    styleUrls: ['./notfound.component.css'],
    imports: [FooterComponent, TranslatePipe],
    standalone: true
})

export class NotFoundComponent {
    constructor(private translate: TranslateService) {
        this.translate.setDefaultLang('en');
        this.translate.use('en');
    }
}