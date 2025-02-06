import { Component } from '@angular/core';
import {RouterLink} from "@angular/router";

@Component({
  selector: 'app-nav',
  templateUrl: './nav.component.html',
  standalone: true,
  imports: [
    RouterLink
  ],
  styleUrl: './nav.component.css'
})
export class NavComponent {}
