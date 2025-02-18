/**********************************/
/* @since 01/01/2025              */
/* @author K10s Open Source Team  */
/**********************************/

import { Component, Input } from '@angular/core';

@Component({
  selector: 'app-loading',
  templateUrl: './loading.component.html',
  standalone: true,
  styleUrls: ['./loading.component.css']
})
export class LoadingComponent {
  @Input() size: number = 24;
}
