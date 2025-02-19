import { Pipe, PipeTransform } from '@angular/core';
import {Node} from "../domain/Kubernetes";

@Pipe({
    standalone: true,
    name: 'nodeCast'
})
export class NodeCastPipe implements PipeTransform {

    transform(value: any, ...args: unknown[]): Node[] {
        return value as Node[];
    }

}
