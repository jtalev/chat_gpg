function togglePageContent() {
    const pageContent = document.getElementById('page-content')
    const navContainer = document.getElementById('nav-container')
    let isPageContentHidden = pageContent.classList.contains('hidden')
    let isNavContainerHidden = navContainer.classList.contains('hidden')

    if (window.innerWidth >= 1030) {
        if (isPageContentHidden) {
            pageContent.classList.remove('hidden')
            pageContent.classList.add('page-content')
        }
        if (isNavContainerHidden) {
            navContainer.classList.remove('hidden')
        }
    } else {
        if (isPageContentHidden) {
            pageContent.classList.remove('hidden')
            pageContent.classList.add('page-content')
        }
        if (!isNavContainerHidden) {
            navContainer.classList.add('hidden')
        }
    }
}

window.addEventListener('resize', togglePageContent)
window.addEventListener('load', togglePageContent)