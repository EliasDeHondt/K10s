/**********************************/
/* @since 01/01/2025              */
/* @author K10s Open Source Team  */
/**********************************/

// Context Menu
let selectedText = '';
let contextMenu = '';
document.addEventListener('DOMContentLoaded', function() {
    document.addEventListener('contextmenu', (event) => {
        if (window.innerWidth < 768) return; // Disable context menu on mobile devices
        event.preventDefault();
        selectedText = window.getSelection().toString(); // Get selected text (For copySelectedText function)
        contextMenu = document.getElementById('context-menu');
        let top = parseInt(contextMenu.style.top);
        let left = parseInt(contextMenu.style.left);

        if (isNaN(top)) top = 0;
        if (isNaN(left)) left = 0;

        if (window.scrollY !== 0) top = event.clientY + window.scrollY;
        else top = event.clientY;

        if (window.scrollX !== 0) left = event.clientX + window.scrollX;
        else left = event.clientX;

        contextMenu.style.top = top + 'px';
        contextMenu.style.left = left + 'px';
        contextMenu.style.display = 'block';

        document.addEventListener('click', (clickEvent) => {
            if (!contextMenu.contains(clickEvent.target)) contextMenu.style.display = 'none';
        });
    });
});

// Copy the current URL to clipboard.
function copyLinkAddress() {
    navigator.clipboard.writeText(window.location.href);
    contextMenu.style.display = 'none';
}

// Copy the selected text to clipboard.
function copySelectedText() {
    if (selectedText) navigator.clipboard.writeText(selectedText);
    contextMenu.style.display = 'none';
}

// Dark Mode
function toggleTheme() {
    const htmlElement = document.documentElement;
    const currentTheme = htmlElement.getAttribute('data-theme');
    const newTheme = currentTheme === 'light' ? 'dark' : 'light';
    htmlElement.setAttribute('data-theme', newTheme);
    localStorage.setItem('theme', newTheme);

}

// Instellen bij pagina laden
document.addEventListener('DOMContentLoaded', () => {
    const savedTheme = localStorage.getItem('theme') || 'light';
    document.documentElement.setAttribute('data-theme', savedTheme);
});