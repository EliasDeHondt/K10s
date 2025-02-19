/**********************************/
/* @since 01/01/2025              */
/* @author K10s Open Source Team  */
/**********************************/

import {Component, OnInit} from '@angular/core';
import {Router, RouterLink} from '@angular/router';
import {CommonModule} from '@angular/common';
import {TranslatePipe, TranslateService} from "@ngx-translate/core";
import {AuthService} from "../services/auth.service";

@Component({
    selector: 'app-nav',
    templateUrl: './nav.component.html',
    standalone: true,
    imports: [RouterLink, CommonModule, TranslatePipe],
    styleUrls: ['./nav.component.css']
})
export class NavComponent implements OnInit {
    githubStars: string = '⭐ Loading...';
    dropdownOpen: boolean = false;
    currentLanguage: string = 'en';

    settingsConfig = {
        languages: [
            {code: 'en', name: 'English'},
            {code: 'nl', name: 'Nederlands'},
            {code: 'fr', name: 'Français'},
            {code: 'de', name: 'Deutsch'},
            {code: 'zh', name: '中文'}
        ]
    };

    constructor(private translate: TranslateService, private authService: AuthService, private router: Router) {
        const savedLang = localStorage.getItem('language');
        if (savedLang) {
            this.translate.use(savedLang);
            this.currentLanguage = savedLang;
        } else {
            this.translate.use('en');
        }
    }

    ngOnInit(): void {
        this.fetchGitHubStars();
    }

    async fetchGitHubStars() {
        try {
            const response = await fetch('https://api.github.com/repos/EliasDeHondt/K10s', {
                headers: {'User-Agent': 'Mozilla/5.0'}
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
        this.translate.use(languageCode);
        localStorage.setItem('language', languageCode);
        this.currentLanguage = languageCode;
        this.closeSettingsModal();
        this.dropdownOpen = false;
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

    logout() {
        this.authService.logout().subscribe({
                next: () => {
                    this.router.navigate(['/login']);
                }
            }
        )
    }

}