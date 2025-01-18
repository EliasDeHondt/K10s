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
    const dashboardHtml = document.documentElement;
    const dashboardMain = document.getElementById('dashboard-main');
    const dashboardTitle = document.getElementById('dashboard-title');

    if (!document.fullscreenElement) {
        if (dashboardHtml.requestFullscreen) dashboardHtml.requestFullscreen();
        else if (dashboardHtml.mozRequestFullScreen) dashboardHtml.mozRequestFullScreen();
        else if (dashboardHtml.webkitRequestFullscreen) dashboardHtml.webkitRequestFullscreen();
        else if (dashboardHtml.msRequestFullscreen) dashboardHtml.msRequestFullscreen();

        dashboardMain.classList.add('fullscreen');
        dashboardMain.style.gridArea = '1 / 1 / -1 / -1';
        dashboardTitle.style.display = 'block';
    } else {
        if (document.exitFullscreen) document.exitFullscreen();
        else if (document.mozCancelFullScreen) document.mozCancelFullScreen();
        else if (document.webkitExitFullscreen) document.webkitExitFullscreen();
        else if (document.msExitFullscreen) document.msExitFullscreen();

        dashboardMain.classList.remove('fullscreen');
        dashboardMain.style.gridArea = '';
        dashboardTitle.style.display = 'none';
    }
});