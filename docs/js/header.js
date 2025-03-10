/**********************************/
/* @since 01/01/2025              */
/* @author K10s Open Source Team  */
/**********************************/

// Fetch GitHub Stars
async function fetchGitHubStars() {
    try {
        const response = await fetch("https://api.github.com/repos/EliasDeHondt/K10s", {
            headers: { "User-Agent": "Mozilla/5.0" }
        });
        if (!response.ok) throw new Error("GitHub API request failed");
        const data = await response.json();
        document.getElementById("header-github-stars").textContent = `⭐ ${data.stargazers_count}`;
    } catch (error) {
        document.getElementById("header-github-stars").textContent = "Error fetching stars";
    }
}