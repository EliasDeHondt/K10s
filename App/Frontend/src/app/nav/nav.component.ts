/**********************************/
/* @since 01/01/2025              */
/* @author K10s Open Source Team  */
/**********************************/

import { Component, OnInit } from '@angular/core';
import { RouterLink } from '@angular/router';
import { CommonModule } from '@angular/common';
import {TranslateService} from "@ngx-translate/core";

@Component({
    selector: 'app-nav',
    templateUrl: './nav.component.html',
    standalone: true,
    imports: [RouterLink, CommonModule],
    styleUrls: ['./nav.component.css']
})
export class NavComponent implements OnInit {
    githubStars: string = '⭐ Loading...';
    dropdownOpen: boolean = false;

    settingsConfig = {
        title: 'Settings',
        languages: [
            { code: 'en', name: 'English' },
            { code: 'nl', name: 'Dutch' },
            { code: 'fr', name: 'French' },
            { code: 'de', name: 'German' },
            { code: 'zh', name: 'Chinese' }
        ]
    };

    constructor(private translate: TranslateService) {
        // Set default language
        this.translate.setDefaultLang('en');

        // Load preferred language from localStorage (if available)
        const savedLang = localStorage.getItem('language');
        if (savedLang) {
            this.translate.use(savedLang);
        }
    }

    ngOnInit(): void {
        this.fetchGitHubStars();
    }

    async fetchGitHubStars() {
        try {
            const response = await fetch('https://api.github.com/repos/EliasDeHondt/K10s', {
                headers: { 'User-Agent': 'Mozilla/5.0' }
            });
            if (!response.ok) throw new Error('GitHub API request failed');

            const data = await response.json();
            this.githubStars = `⭐ ${data.stargazers_count}`;
        } catch (error) {
            this.githubStars = '❌ Error fetching stars';
        }
    }

    openSettingsModal() {
        document.querySelector('.modal')?.classList.add('show');
    }

    closeSettingsModal() {
        document.querySelector('.modal')?.classList.remove('show');
    }

    changeLanguage(languageCode: string) {
        // Change the language based on the selected option
        this.translate.use(languageCode);
        localStorage.setItem('language', languageCode); // Save the selected language to localStorage
        this.closeSettingsModal(); // Close the modal after language change
    }

    toggleDropdown() {
        this.dropdownOpen = !this.dropdownOpen;
    }

    toggleTheme() {
        const htmlElement = document.documentElement;
        const currentTheme = htmlElement.getAttribute('data-theme');
        const newTheme = currentTheme === 'light' ? 'dark' : 'light';
        htmlElement.setAttribute('data-theme', newTheme);
        localStorage.setItem('theme', newTheme);
    }
}
