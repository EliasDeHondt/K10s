import { Component } from '@angular/core';
import {NavComponent} from '../nav/nav.component';

@Component({
  selector: 'app-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.css'],
  imports: [NavComponent],
  standalone: true
})
export class DashboardComponent {
  title = 'K10s';
}
