import { bootstrapApplication } from '@angular/platform-browser';
import { appConfig } from './app/app.config';
import {DashboardComponent} from "./app/components/dashboard/dashboard.component";

bootstrapApplication(DashboardComponent, appConfig)
    .catch((err) => console.error(err));