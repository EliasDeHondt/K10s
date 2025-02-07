/**********************************/
/* @since 01/01/2025              */
/* @author K10s Open Source Team  */
/**********************************/

import { provideRouter } from '@angular/router';
import { routes } from './app.routes';
import { Title } from '@angular/platform-browser';

export const appConfig = {
    providers: [
        provideRouter(routes),
        { provide: Title, useValue: 'K10s' }
    ]
};