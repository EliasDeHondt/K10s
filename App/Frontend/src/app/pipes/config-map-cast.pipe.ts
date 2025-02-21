import { Pipe, PipeTransform } from '@angular/core';
import {ConfigMap} from "../domain/Kubernetes";

@Pipe({
  standalone: true,
  name: 'configMapCast'
})
export class ConfigMapCastPipe implements PipeTransform {

  transform(value: unknown, ...args: unknown[]): ConfigMap[] {
    return value as ConfigMap[];
  }

}
