import { Pipe, PipeTransform } from '@angular/core';
import {Service} from "../domain/Kubernetes";

@Pipe({
  standalone: true,
  name: 'serviceCast'
})
export class ServiceCastPipe implements PipeTransform {

  transform(value: unknown, ...args: unknown[]): Service[] {
    return value as Service[];
  }

}
