/**********************************/
/* @since 01/01/2025              */
/* @author K10s Open Source Team  */
/**********************************/

// Load external content
document.addEventListener('DOMContentLoaded', function() {
    loadExternalContent("nav", "/includes/nav.html");
    loadExternalContent("footer", "/includes/footer.html");
    loadExternalContent("context-menu", "/includes/context-menu.html");
});

// Fullscreen button
document.getElementById('dashboard-fullscreen-button').addEventListener('click', () => {
    const dashboardMain = document.getElementById('dashboard-main');
    const dashboardTitle = document.getElementById('dashboard-title');
    if (!dashboardMain.classList.contains('fullscreen')) {
        dashboardMain.classList.add('fullscreen');
        dashboardMain.style.gridArea = '1 / 1 / -1 / -1';
        dashboardTitle.style.display = 'block';
    } else {
        dashboardMain.classList.remove('fullscreen');
        dashboardMain.style.gridArea = '';
        dashboardTitle.style.display = 'none';
    }
});