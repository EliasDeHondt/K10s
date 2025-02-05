import anime from 'animejs';
import {Component, ViewEncapsulation} from '@angular/core';
import { FormsModule } from '@angular/forms';
import {loadExternalContent} from '../../main';

@Component({
  selector: 'app-login',
  standalone: true,
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.css'],
  imports: [FormsModule],
})
export class LoginComponent {
  username: string = '';
  password: string = '';

  onSubmit() {
    console.log('Login attempt:', { username: this.username, password: this.password });
  }
}


//todo: elias code

// // Load external content
// document.addEventListener('DOMContentLoaded', function() {
//   loadExternalContent("context-menu", "/includes/context-menu.html");
// });

// Background animations for login page.
document.querySelectorAll<SVGElement>('g').forEach(function(cube, index) {
  const transform = cube.getAttribute('transform');
  let currentTranslateY = 0;
  let currentTranslateX = 0;
  let scale = 1;

  if (transform) {
    const transformValues = transform.split('(')[1].split(')')[0].split(',');

    currentTranslateX = parseFloat(transformValues[0]) || 0;
    currentTranslateY = parseFloat(transformValues[1]) || 0;

    const scaleValue = transform.split('scale(')[1]?.split(')')[0];
    scale = scaleValue ? parseFloat(scaleValue) : 1;
  }

  anime({
    targets: cube,
    translateY: [currentTranslateY, currentTranslateY - 150],
    translateX: [currentTranslateX, currentTranslateX], // No animation
    scale: [scale, scale], // No animation
    duration: 1500, // 1.5 seconds
    direction: 'alternate',
    loop: true,
    delay: index * 100,
    endDelay: (el, i, l) => (l - i) * 100
  });
});
