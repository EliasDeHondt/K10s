import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  standalone: true,
  name: 'byteFormat'
})
export class ByteFormatPipe implements PipeTransform {
  transform(value: number): string {
    if (!value) {
      return '0 GB';
    }

    const GB = 1024 * 1024 * 1024;
    const TB = GB * 1000;

    if (value >= TB) {
      return (value / TB).toFixed(2) + ' TB';
    } else {
      return (value / GB).toFixed(2) + ' GB';
    }
  }
}
