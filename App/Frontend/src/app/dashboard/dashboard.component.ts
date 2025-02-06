import { Component,ViewChild } from '@angular/core';
import {NavComponent} from '../nav/nav.component';
import {FooterComponent} from "../footer/footer.component";
import {ContextMenuComponent} from "../context-menu/context-menu.component";

@Component({
  selector: 'app-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.css'],
  imports: [NavComponent, FooterComponent, ContextMenuComponent],
  standalone: true
})
export class DashboardComponent {
  title = 'K10s';
  @ViewChild(ContextMenuComponent) contextMenu!: ContextMenuComponent;

  onRightClick(event: MouseEvent) {
    event.preventDefault(); // Stops the default context menu
    this.contextMenu?.show(event);
  }

}
