/**********************************/
/* @since 01/01/2025              */
/* @author K10s Open Source Team  */
/**********************************/

import { Routes } from '@angular/router';
import { LoginComponent } from './login/login.component';
import { DashboardComponent } from './dashboard/dashboard.component';
import { SearchComponent } from './search/search.component';
import { AddComponent } from "./add/add.component";
import { NotFoundComponent } from './notfound/notfound.component';

export const routes: Routes = [
    { path: '', redirectTo: 'login', pathMatch: 'full' },
    { path: 'login', component: LoginComponent, data: { title: 'K10s - Login' } },
    { path: 'dashboard', component: DashboardComponent, data: { title: 'K10s - Dashboard' } },
    { path: 'search', component: SearchComponent, data: { title: 'K10s - Search' } },
    { path: 'add', component: AddComponent, data: { title: 'K10s - Add' } },
    { path: '**', component: NotFoundComponent, data: { title: 'K10s - Not Found' } }
];