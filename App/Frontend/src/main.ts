/**********************************/
/* @since 01/01/2025              */
/* @author K10s Open Source Team  */
/**********************************/

import { bootstrapApplication } from '@angular/platform-browser';
import { appConfig } from './app/app.config';
import {LoginComponent} from "./app/login/login.component";

bootstrapApplication(LoginComponent, appConfig)
    .catch((err) => console.error(err));

const htmlElement = document.documentElement;
htmlElement.setAttribute('data-theme', 'light');



// // Context Menu
// let selectedText: string = '';
// let contextMenu: HTMLElement | null = null;
//
// document.addEventListener('DOMContentLoaded', function() {
//   document.addEventListener('contextmenu', (event: MouseEvent) => {
//     if (window.innerWidth < 768) return; // Disable context menu on mobile devices
//     event.preventDefault();
//     selectedText = window.getSelection()?.toString() || ''; // Get selected text (For copySelectedText function)
//     contextMenu = document.getElementById('context-menu');
//
//     if (contextMenu) {
//       let top: number = parseInt(contextMenu.style.top);
//       let left: number = parseInt(contextMenu.style.left);
//
//       if (isNaN(top)) top = 0;
//       if (isNaN(left)) left = 0;
//
//       if (window.scrollY !== 0) top = event.clientY + window.scrollY;
//       else top = event.clientY;
//
//       if (window.scrollX !== 0) left = event.clientX + window.scrollX;
//       else left = event.clientX;
//
//       contextMenu.style.top = `${top}px`;
//       contextMenu.style.left = `${left}px`;
//       contextMenu.style.display = 'block';
//
//       document.addEventListener('click', (clickEvent: MouseEvent) => {
//         if (contextMenu && !contextMenu.contains(clickEvent.target as Node)) {
//           contextMenu.style.display = 'none';
//         }
//       });
//     }
//   });
// });
//
// // Copy the current URL to clipboard.
// function copyLinkAddress(): void {
//   navigator.clipboard.writeText(window.location.href);
//   if (contextMenu) contextMenu.style.display = 'none';
// }
//
// // Copy the selected text to clipboard.
// function copySelectedText(): void {
//   if (selectedText) navigator.clipboard.writeText(selectedText);
//   if (contextMenu) contextMenu.style.display = 'none';
// }
//
// // Dark Mode
// function toggleTheme(): void {
//   const htmlElement: HTMLElement = document.documentElement;
//   const currentTheme: string | null = htmlElement.getAttribute('data-theme');
//   const newTheme: string = currentTheme === 'light' ? 'dark' : 'light';
//   htmlElement.setAttribute('data-theme', newTheme);
//   localStorage.setItem('theme', newTheme);
// }
//
// // Set on page load
// document.addEventListener('DOMContentLoaded', () => {
//   const savedTheme: string | null = localStorage.getItem('theme') || 'light';
//   document.documentElement.setAttribute('data-theme', savedTheme);
// });
