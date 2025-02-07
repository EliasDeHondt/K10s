/**********************************/
/* @since 01/01/2025              */
/* @author K10s Open Source Team  */
/**********************************/

import { bootstrapApplication } from '@angular/platform-browser';
import { AppComponent } from './app/app.component';
import { appConfig } from './app/app.config';

bootstrapApplication(AppComponent, appConfig)
    .catch(err => console.error(err));

const htmlElement = document.documentElement;
htmlElement.setAttribute('data-theme', 'light');

// Content Loader
export function loadExternalContent(DivId: string, url: string): void {
    let xmlhttp: XMLHttpRequest | ActiveXObject;
    if (window.XMLHttpRequest) xmlhttp = new XMLHttpRequest();
    else xmlhttp = new ActiveXObject("Microsoft.XMLHTTP");

    // Type assertion to tell TypeScript that xmlhttp is an XMLHttpRequest here
    (xmlhttp as XMLHttpRequest).onreadystatechange = function(): void {
        if ((xmlhttp as XMLHttpRequest).readyState === XMLHttpRequest.DONE) {
        if ((xmlhttp as XMLHttpRequest).status === 200) {
            const div = document.getElementById(DivId);
            if (div) {
            div.innerHTML = (xmlhttp as XMLHttpRequest).responseText;
            let scripts = div.getElementsByTagName('script');
            for (let i = 0; i < scripts.length; i++) {
                let script = document.createElement('script');
                script.text = scripts[i].text;
                document.body.appendChild(script);
            }
            }
        }
        }
    };

    (xmlhttp as XMLHttpRequest).open("GET", url, true);
    (xmlhttp as XMLHttpRequest).send();
}

// Open Modal
function openModal(Id: string, modalHTML: string): void {
    document.body.insertAdjacentHTML('beforeend', modalHTML);
    const modal = document.getElementById(Id);
    if (modal) modal.style.display = 'block';
}

// Close Modal
function closeModal(Id: string): void {
    const modal = document.getElementById(Id);
    if (modal) modal.remove();
}

// Toggle Dropdown
function toggleDropdown(id: string): void {
    const dropdown = document.querySelector(id);
    if (dropdown) dropdown.classList.toggle('show');
}

// Change Language
function changeLanguage(language: string): void {
    // TODO
    translatePage(language);
}

// Translate Page
function translatePage(language: string): void {
    // TODO
}

// Fetch GitHub Stars
async function fetchGitHubStars(): Promise<void> {
    try {
        const response = await fetch("https://api.github.com/repos/EliasDeHondt/K10s", {
            headers: { "User-Agent": "Mozilla/5.0" }
        });
        if (!response.ok) throw new Error("GitHub API request failed");
        const data = await response.json();
        const starsElement = document.getElementById("nav-github-stars");
        if (starsElement) starsElement.textContent = `⭐ ${data.stargazers_count}`;
    } catch (error) {
        const starsElement = document.getElementById("nav-github-stars");
        if (starsElement) starsElement.textContent = "Error fetching stars";
    }
}

// Context Menu
let selectedText: string = '';
let contextMenu: HTMLElement | null = null;

document.addEventListener('DOMContentLoaded', function() {
    document.addEventListener('contextmenu', (event: MouseEvent) => {
        if (window.innerWidth < 768) return; // Disable context menu on mobile devices
        event.preventDefault();
        selectedText = window.getSelection()?.toString() || ''; // Get selected text (For copySelectedText function)
        contextMenu = document.getElementById('context-menu');

        if (contextMenu) {
        let top: number = parseInt(contextMenu.style.top);
        let left: number = parseInt(contextMenu.style.left);

        if (isNaN(top)) top = 0;
        if (isNaN(left)) left = 0;

        if (window.scrollY !== 0) top = event.clientY + window.scrollY;
        else top = event.clientY;

        if (window.scrollX !== 0) left = event.clientX + window.scrollX;
        else left = event.clientX;

        contextMenu.style.top = `${top}px`;
        contextMenu.style.left = `${left}px`;
        contextMenu.style.display = 'block';

        document.addEventListener('click', (clickEvent: MouseEvent) => {
            if (contextMenu && !contextMenu.contains(clickEvent.target as Node)) {
            contextMenu.style.display = 'none';
            }
        });
        }
    });
});

// Copy the current URL to clipboard.
function copyLinkAddress(): void {
    navigator.clipboard.writeText(window.location.href);
    if (contextMenu) contextMenu.style.display = 'none';
}

// Copy the selected text to clipboard.
function copySelectedText(): void {
    if (selectedText) navigator.clipboard.writeText(selectedText);
    if (contextMenu) contextMenu.style.display = 'none';
}

// Dark Mode
function toggleTheme(): void {
    const htmlElement: HTMLElement = document.documentElement;
    const currentTheme: string | null = htmlElement.getAttribute('data-theme');
    const newTheme: string = currentTheme === 'light' ? 'dark' : 'light';
    htmlElement.setAttribute('data-theme', newTheme);
    localStorage.setItem('theme', newTheme);
}

// Set on page load
document.addEventListener('DOMContentLoaded', () => {
    const savedTheme: string | null = localStorage.getItem('theme') || 'light';
    document.documentElement.setAttribute('data-theme', savedTheme);
});