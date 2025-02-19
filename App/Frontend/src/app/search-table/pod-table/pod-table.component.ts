import {Component, Input} from '@angular/core';
import {TranslatePipe} from "@ngx-translate/core";
import {Pod} from "../../domain/Kubernetes";

@Component({
  selector: 'app-pod-search-table',
  imports: [
    TranslatePipe
  ],
  templateUrl: './pod-table.component.html',
  standalone: true,
  styleUrls: ['./pod-table.component.css', '../../search/search.component.css']
})
export class PodTableComponent {

  @Input({required: true}) pods!: Pod[];

}
