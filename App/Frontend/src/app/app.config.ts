import { provideRouter } from '@angular/router';
import { routes } from './app.routes';

export const appConfig = {
    providers: [
        provideRouter(routes) // ✅ Ensure routes are provided
    ]
};
