/**********************************/
/* @since 01/01/2025              */
/* @author K10s Open Source Team  */
/**********************************/

// Load external content
document.addEventListener('DOMContentLoaded', function() {
    loadExternalContent("nav", "/includes/nav.html");
    loadExternalContent("footer", "/includes/footer.html");
    loadExternalContent("context-menu", "/includes/context-menu.html");
    fetchGitHubStars();
});