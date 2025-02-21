import { Pipe, PipeTransform } from '@angular/core';
import {Secret} from "../domain/Kubernetes";

@Pipe({
  standalone: true,
  name: 'secretCast'
})
export class SecretCastPipe implements PipeTransform {

  transform(value: unknown, ...args: unknown[]): Secret[] {
    return value as Secret[];
  }

}
