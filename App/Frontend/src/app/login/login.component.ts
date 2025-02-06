import anime from 'animejs';
import {Component} from '@angular/core';
import { FormsModule } from '@angular/forms';
import {Router, RouterModule} from "@angular/router";
import {FooterComponent} from "../footer/footer.component";
// import {loadExternalContent} from '../../main';

@Component({
  selector: 'app-login',
  standalone: true,
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.css'],
  imports: [FormsModule, RouterModule, FooterComponent],
})
export class LoginComponent {
  username: string = '';
  password: string = '';
  constructor(private router: Router) {}

  onSubmit() { //todo met backend
    console.log(this.username, this.password);
    if (this.username && this.password) {
      this.router.navigate(['/dashboard']);
    } else {
      alert('Please enter valid credentials.');
    }
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
