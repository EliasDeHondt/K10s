/**********************************/
/* @since 01/01/2025              */
/* @author K10s Open Source Team  */
/**********************************/

// Settings configuration
const settingsConfig = {
    title: "Settings",
    languages: [
        { code: 'en', name: 'English' },
        { code: 'nl', name: 'Dutch' },
        { code: 'fr', name: 'French' },
        { code: 'de', name: 'German' },
        { code: 'zh', name: 'Chinese' }
    ]
};

// Function to generate the settings modal HTML
function generateSettingsModalHTML(config) {
    const languageOptions = config.languages.map(lang => `
        <li class="cursor-pointer" onclick="changeLanguage('${lang.code}'); toggleDropdown('.modal-dropdown')">${lang.name}</li>
    `).join('');

    return `
        <section id="settingsModal" class="modal">
            <article class="modal-article">
                <span class="modal-close-button cursor-pointer" onclick="closeModal('settingsModal')">x</span>
                <h2>${config.title}</h2>
                <br>
                <button class="modal-theme-toggle-button cursor-pointer" onclick="toggleTheme()">
                    <svg viewBox="0 0 24 24" class="cursor-pointer">
                        <path d="m12 20.501v-17.001c-.019 0-.041 0-.064 0-1.549 0-2.999.424-4.24 1.162l.038-.021c-2.554 1.5-4.242 4.233-4.242 7.36s1.688 5.86 4.202 7.339l.04.022c1.202.716 2.652 1.14 4.2 1.14h.07-.004zm12.001-8.501v.09c0 2.187-.598 4.235-1.64 5.988l.03-.054c-1.067 1.824-2.544 3.301-4.311 4.337l-.056.031c-1.729 1.012-3.807 1.609-6.024 1.609s-4.295-.597-6.081-1.64l.057.031c-1.824-1.067-3.301-2.544-4.337-4.311l-.031-.056c-1.012-1.729-1.609-3.807-1.609-6.024s.597-4.295 1.64-6.081l-.031.057c1.067-1.824 2.544-3.301 4.311-4.337l.056-.031c1.729-1.012 3.807-1.609 6.024-1.609s4.295.597 6.081 1.64l-.057-.031c1.824 1.067 3.301 2.544 4.337 4.311l.031.056c1.012 1.699 1.61 3.747 1.61 5.934v.095z">
                    </svg>
                </button>
                <br><br>
                <section class="modal-dropdown">
                    <button class="modal-dropdown-button cursor-pointer" onclick="toggleDropdown('.modal-dropdown')">Choose Language</button>
                    <ul class="modal-dropdown-menu cursor-pointer">
                        ${languageOptions}
                    </ul>
                </section>
            </article>
        </section>
    `;
}

// Generate the settings modal HTML
const settingsModalHTML = generateSettingsModalHTML(settingsConfig);