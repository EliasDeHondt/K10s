/**********************************/
/* @since 01/01/2025              */
/* @author K10s Open Source Team  */
/**********************************/

import { bootstrapApplication } from '@angular/platform-browser';
import { AppComponent } from './app/app.component';
import { appConfig } from './app/app.config';
import { TranslateHttpLoader } from "@ngx-translate/http-loader";
import { HttpClient } from "@angular/common/http";
import { TranslateLoader } from "@ngx-translate/core";

export function HttpLoaderFactory(http: HttpClient): TranslateLoader {
    return new TranslateHttpLoader(http, 'assets/i18n/', '.json');
}

bootstrapApplication(AppComponent, appConfig).catch(err => console.error(err));

document.documentElement.setAttribute('data-theme', 'light');

// Content Loader
export function loadExternalContent(DivId: string, url: string): void {
    let xmlhttp: XMLHttpRequest | ActiveXObject;
    if (window.XMLHttpRequest) xmlhttp = new XMLHttpRequest();
    else xmlhttp = new ActiveXObject("Microsoft.XMLHTTP");

    (xmlhttp as XMLHttpRequest).onreadystatechange = function(): void {
        if ((xmlhttp as XMLHttpRequest).readyState === XMLHttpRequest.DONE) {
            if ((xmlhttp as XMLHttpRequest).status === 200) {
                const div = document.getElementById(DivId);
                if (div) {
                    div.innerHTML = (xmlhttp as XMLHttpRequest).responseText;
                    let scripts = div.getElementsByTagName('script');
                    for (let i = 0; i < scripts.length; i++) {
                        let script = document.createElement('script');
                        script.text = scripts[i].text;
                        document.body.appendChild(script);
                    }
                }
            }
        }
    };

    (xmlhttp as XMLHttpRequest).open("GET", url, true);
    (xmlhttp as XMLHttpRequest).send();
}

// Set on page load
document.addEventListener('DOMContentLoaded', () => {
    const savedTheme: string | null = localStorage.getItem('theme') || 'light';
    document.documentElement.setAttribute('data-theme', savedTheme);
});