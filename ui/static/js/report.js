function onLeaveRequestClick() {
	const requestRow = document.querySelector("#report-leave-request-row-unclicked")
	const expandedRow = document.querySelector("#report-leave-request-row-expanded")
	const computedStyle = window.getComputedStyle(requestRow)
	if (computedStyle.display === "flex") {
		requestRow.style.display = "none"
		expandedRow.style.display = "flex"
	} else if (computedStyle.display === "none") {
		requestRow.style.display = "flex"
		expandedRow.style.display = "none"
	}
}