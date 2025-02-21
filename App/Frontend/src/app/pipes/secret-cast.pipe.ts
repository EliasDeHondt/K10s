/**********************************/
/* @since 01/01/2025              */
/* @author K10s Open Source Team  */
/**********************************/

import { Pipe, PipeTransform } from '@angular/core';
import { Secret } from "../domain/Kubernetes";

@Pipe({
    standalone: true,
    name: 'secretCast'
})

export class SecretCastPipe implements PipeTransform {
    transform(value: unknown, ...args: unknown[]): Secret[] {
        return value as Secret[];
    }
}