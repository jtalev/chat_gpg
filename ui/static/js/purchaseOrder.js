function toggleContentSelector(event) {
	const contentSelectors = document.querySelectorAll(".view-content-selector")
	contentSelectors.forEach(selector => {
		selector.style.backgroundColor = "var(--main-background-color)"
		selector.style.borderBottom = "solid 1px var(--gpg-green)"
	})

	const target = event.currentTarget
	target.style.backgroundColor = "white"
	target.style.borderBottom = "none"
}