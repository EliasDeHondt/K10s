import { Component, ElementRef, HostListener } from '@angular/core';
import {NgIf} from '@angular/common';

@Component({
  selector: 'app-context-menu',
  standalone: true,
  templateUrl: './context-menu.component.html',
  styleUrls: ['./context-menu.component.css'],
  imports: [
    NgIf
  ]
})
export class ContextMenuComponent {
  isVisible = false;
  x = 0;
  y = 0;

  constructor(private elementRef: ElementRef) {}

  show(event: MouseEvent) {
    event.preventDefault();
    this.isVisible = true;
    this.x = event.clientX;
    this.y = event.clientY;
  }

  hide() {
    this.isVisible = false;
  }

  @HostListener('document:click', ['$event'])
  onClickOutside(event: Event) {
    if (!this.elementRef.nativeElement.contains(event.target)) {
      this.hide();
    }
  }
}
