import { Pipe, PipeTransform } from '@angular/core';
import {Deployment} from "../domain/Kubernetes";

@Pipe({
  standalone: true,
  name: 'deploymentCast'
})
export class DeploymentCastPipe implements PipeTransform {

  transform(value: unknown, ...args: unknown[]): Deployment[] {
    return value as Deployment[];
  }

}
