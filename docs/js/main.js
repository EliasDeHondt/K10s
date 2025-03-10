/**********************************/
/* @since 01/01/2025              */
/* @author K10s Open Source Team  */
/**********************************/

// Content Loader
function loadExternalContent(DivId, url) {
    let xmlhttp;
    if (window.XMLHttpRequest) xmlhttp = new XMLHttpRequest();
    else xmlhttp = new ActiveXObject("Microsoft.XMLHTTP");
    xmlhttp.onreadystatechange = function() {
        if (xmlhttp.readyState == XMLHttpRequest.DONE ) {
            if(xmlhttp.status == 200) {
                document.getElementById(DivId).innerHTML = xmlhttp.responseText;
                let scripts = document.getElementById(DivId).getElementsByTagName('script');
                for (let i = 0; i < scripts.length; i++) {
                    let script = document.createElement('script');
                    script.text = scripts[i].text;
                    document.body.appendChild(script);
                }
            }
        }
    }
    xmlhttp.open("GET", url, true);
    xmlhttp.send();
}

// Load external content
document.addEventListener('DOMContentLoaded', function() {
    loadExternalContent("header", "/includes/header.html");
    loadExternalContent("footer", "/includes/footer.html");
    loadExternalContent("context-menu", "/includes/context-menu.html");
    fetchGitHubStars();
});