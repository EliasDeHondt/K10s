import { Component, AfterViewInit } from '@angular/core';
import { Router, RouterModule } from '@angular/router';
import { FormsModule } from '@angular/forms';
import anime from 'animejs/lib/anime.es.js';
import { FooterComponent } from '../footer/footer.component';

@Component({
  selector: 'app-login',
  standalone: true,
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.css'],
  imports: [FormsModule, RouterModule, FooterComponent],
})
export class LoginComponent implements AfterViewInit {
  username: string = '';
  password: string = '';

  constructor(private router: Router) {}

  ngAfterViewInit() {
    const cubes = document.querySelectorAll('g');
    cubes.forEach((cube, index) => {
      const transform = cube.getAttribute('transform') || 'translate(0,0) scale(1)';
      const translateMatch = transform.match(/translate\(([^,]+),([^,]+)\)/);
      const scaleMatch = transform.match(/scale\(([^)]+)\)/);

      const currentTranslateX = translateMatch ? parseFloat(translateMatch[1]) : 0;
      const currentTranslateY = translateMatch ? parseFloat(translateMatch[2]) : 0;
      const scale = scaleMatch ? parseFloat(scaleMatch[1]) : 1;

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
  }

  onSubmit() {
    console.log(this.username, this.password);
    if (this.username && this.password) {
      this.router.navigate(['/dashboard']);
    } else {
      alert('Please enter valid credentials.');
    }
  }
}
