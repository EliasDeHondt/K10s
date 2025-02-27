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
import { environment } from './environments/environment';

export function HttpLoaderFactory(http: HttpClient): TranslateLoader {
    return new TranslateHttpLoader(http, 'assets/i18n/', '.json');
}

bootstrapApplication(AppComponent, appConfig).catch(err => console.error(err));

document.documentElement.setAttribute('data-theme', 'light');

// Set on page load
document.addEventListener('DOMContentLoaded', () => {
    const savedTheme: string | null = localStorage.getItem('theme') || 'light';
    document.documentElement.setAttribute('data-theme', savedTheme);
});

if (environment.BASE_URL === 'http://localhost:8080') {
    console.log('You are running in development mode');
} else if (environment.BASE_URL === '') {
    environment.BASE_URL = `${window.location.origin}/api`;
    window.location.reload(); // Refresh the page to apply the new BASE_URL
}