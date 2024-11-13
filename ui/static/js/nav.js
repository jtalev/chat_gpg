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