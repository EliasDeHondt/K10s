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

// Fetch GitHub Stars
async function fetchGitHubStars() {
    try {
        const response = await fetch("https://api.github.com/repos/EliasDeHondt/K10s", {
            headers: { "User-Agent": "Mozilla/5.0" }
        });
        if (!response.ok) throw new Error("GitHub API request failed");
        const data = await response.json();
        document.getElementById("index-header-github-stars").textContent = `⭐ ${data.stargazers_count}`;
    } catch (error) {
        document.getElementById("index-header-github-stars").textContent = "Error fetching stars";
    }
}

// Load external content
document.addEventListener('DOMContentLoaded', function() {
    loadExternalContent("footer", "/includes/footer.html");
    loadExternalContent("context-menu", "/includes/context-menu.html");
    fetchGitHubStars();
});