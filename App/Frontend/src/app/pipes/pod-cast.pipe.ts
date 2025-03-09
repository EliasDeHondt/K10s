/**********************************/
/* @since 01/01/2025              */
/* @author K10s Open Source Team  */
/**********************************/

import { Pipe, PipeTransform } from '@angular/core';
import { Pod } from "../domain/Kubernetes";

@Pipe({
    standalone: true,
    name: 'podCast'
})

export class PodCastPipe implements PipeTransform {
    transform(value: any, ...args: unknown[]): Pod[] {
        return value as Pod[];
    }
}