import {Component, OnInit} from '@angular/core';
import {RouterLink} from "@angular/router";

@Component({
  selector: 'app-nav',
  templateUrl: './nav.component.html',
  standalone: true,
  imports: [
    RouterLink
  ],
  styleUrl: './nav.component.css'
})
export class NavComponent implements OnInit{



  // load github stars
  githubStars: string = "⭐ Loading...";

  constructor() {}

  ngOnInit(): void {
    this.fetchGitHubStars();
  }

  async fetchGitHubStars() {
    try {
      const response = await fetch("https://api.github.com/repos/EliasDeHondt/K10s", {
        headers: { "User-Agent": "Mozilla/5.0" }
      });
      if (!response.ok) throw new Error("GitHub API request failed");

      const data = await response.json();
      this.githubStars = `⭐ ${data.stargazers_count}`;
    } catch (error) {
      this.githubStars = "❌ Error fetching stars";
    }
  }
}
