import { Component } from '@angular/core';
import {NavComponent} from '../nav/nav.component';
import {FooterComponent} from "../footer/footer.component";

@Component({
  selector: 'app-search',
  templateUrl: './search.component.html',
  styleUrls: ['./search.component.css'],
  imports: [NavComponent,FooterComponent],
  standalone: true
})
export class SearchComponent {
  title = 'K10s';
}
