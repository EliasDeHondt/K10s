/**********************************/
/* @since 01/01/2025              */
/* @author K10s Open Source Team  */
/**********************************/

import { Pipe, PipeTransform } from '@angular/core';
import { ConfigMap } from "../domain/Kubernetes";

@Pipe({
    standalone: true,
    name: 'configMapCast'
})

export class ConfigMapCastPipe implements PipeTransform {
    transform(value: unknown, ...args: unknown[]): ConfigMap[] {
        return value as ConfigMap[];
    }
}