function toggleMenu() {
    const pageContent = document.getElementById('page-content')
    const navContainer = document.getElementById('nav-container')
    const isPageContentHidden = pageContent.classList.contains('hidden')

    if (isPageContentHidden) {
        pageContent.classList.remove('hidden')
        pageContent.classList.add('page-content')
        navContainer.classList.add('hidden')
    } else {
        pageContent.classList.add('hidden')
        pageContent.classList.remove('page-content')
        navContainer.classList.remove('hidden')
    }
}

document.addEventListener("DOMContentLoaded", onNavLinkClick) 
function onNavLinkClick() {
    const navLinks = document.querySelectorAll(".nav-link")
    const currentUrl = new URL(window.location.href).pathname;

    navLinks.forEach(element => {
        if ("/" + element.firstChild.textContent === currentUrl) {
            element.style.borderLeft = "var(--gpg-green) solid 3px"
        } else {
            element.style.borderLeft = "var(--gpg-grey) solid 3px"
        }
    })
}