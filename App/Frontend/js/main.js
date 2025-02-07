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

// Open Modal
function openModel(Id, modalHTML) {
    document.body.insertAdjacentHTML('beforeend', modalHTML);
    document.getElementById(Id).style.display = 'block';
}

// Close Modal
function closeModal(Id) {
    const modal = document.getElementById(Id);
    if (modal) modal.remove();
}

// Toggle Dropdown
function toggleDropdown(id) {
    document.querySelector(id).classList.toggle('show');
}

// Change Language
function changeLanguage(language) {
    // TODO
    translatePage(language);
}

// Translate Page
function translatePage(language) {
    // TODO
}
