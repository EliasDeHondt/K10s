/**********************************/
/* @since 01/01/2025              */
/* @author K10s Open Source Team  */
/**********************************/

// Set the current year in the footer
function setTheCurrentYear() {
    var footerYearElement = document.getElementById("footer-year");
    if (footerYearElement) footerYearElement.innerHTML = new Date().getFullYear();
}

setTheCurrentYear();