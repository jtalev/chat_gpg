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

    const menuIcon = document.querySelector(".menuIconArrow")
    if (menuIcon) {
        menuIcon.classList.toggle('rotated')
    }
}

document.addEventListener("DOMContentLoaded", onNavLinkClick) 
function onNavLinkClick() {
    const navLinks = document.querySelectorAll(".nav-link")
    const navLinksText = document.querySelectorAll(".nav-link-text")
    const currentUrl = new URL(window.location.href).pathname;
    const windowWidth = window.innerWidth

    if (windowWidth < 1030) {
        navLinksText.forEach(element => {
            if ("/" + element.firstChild.textContent === currentUrl) {
                element.style.color = "var(--gpg-green)"
            } else {
                element.style.color = "var(--gpg-grey)"
            }
        })
    } else {
        navLinks.forEach(element => {
            if ("/" + element.firstChild.textContent === currentUrl) {
                element.style.borderLeft = "var(--gpg-green) solid 3px"
                element.style.backgroundColor = "#00969C50"
            } else {
                element.style.borderLeft = "var(--gpg-grey) solid 3px"
            }
        })
    }
}