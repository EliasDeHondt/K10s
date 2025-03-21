import { Component } from '@angular/core';
import {NgForOf} from "@angular/common";
import {LoadingService} from "../../services/loading.service";

@Component({
  selector: 'app-table-skeleton',
  imports: [
    NgForOf
  ],
  templateUrl: './table-skeleton.component.html',
  standalone: true,
  styleUrl: '../../search/search.component.css'
})
export class TableSkeletonComponent {

  constructor(protected loadingService: LoadingService){}

}
