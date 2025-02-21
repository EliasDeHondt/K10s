import {Component, Input} from '@angular/core';
import {TranslatePipe} from "@ngx-translate/core";
import {Deployment} from "../../domain/Kubernetes";

@Component({
  selector: 'app-deployment-table',
  imports: [
      TranslatePipe
  ],
  templateUrl: './deployment-table.component.html',
  standalone: true,
  styleUrls: ['./deployment-table.component.css', '../../search/search.component.css']
})
export class DeploymentTableComponent {

  @Input({required: true}) deployments!: Deployment[];

}
