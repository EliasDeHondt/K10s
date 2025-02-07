import { Component, OnInit } from '@angular/core';
import { RouterLink } from "@angular/router";

@Component({
  selector: 'app-nav',
  templateUrl: './nav.component.html',
  standalone: true,
  imports: [
    RouterLink
  ],
  styleUrls: ['./nav.component.css']
})
export class NavComponent implements OnInit {
  githubStars: string = "⭐ Loading...";
  showSettingsModal: boolean = false; // Controls modal visibility
  settingsConfig = {
    title: "Settings",
    languages: [
      { code: 'en', name: 'English' },
      { code: 'nl', name: 'Dutch' },
      { code: 'fr', name: 'French' },
      { code: 'de', name: 'German' },
      { code: 'zh', name: 'Chinese' }
    ]
  };

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

  openSettingsModal() {
    this.showSettingsModal = true;
  }

  closeSettingsModal() {
    this.showSettingsModal = false;
  }

  changeLanguage(languageCode: string) {
    console.log("Language changed to:", languageCode);
    this.closeSettingsModal();
  }

  // Toggle Dropdown
  toggleDropdown(id: string) {
    document.querySelector(id)?.classList.toggle('show');
  }

  // Dark Mode
  toggleTheme() {
    const htmlElement = document.documentElement;
    const currentTheme = htmlElement.getAttribute('data-theme');
    const newTheme = currentTheme === 'light' ? 'dark' : 'light';
    htmlElement.setAttribute('data-theme', newTheme);
    localStorage.setItem('theme', newTheme);
  }
}
