function onLeaveRequestClick(event) {
  const clickedRow = event.currentTarget;

  const requestRow = clickedRow.querySelector(".report-leave-request-row-unclicked");
  const expandedRow = clickedRow.querySelector(".report-leave-request-row-expanded");

  const computedStyle = window.getComputedStyle(requestRow);
  if (computedStyle.display === "flex") {
    requestRow.style.display = "none";
    expandedRow.style.display = "flex";
  } else {
    requestRow.style.display = "flex";
    expandedRow.style.display = "none";
  }
}

function toggleContentSelector(event) {
  const contentSelectors = document.querySelectorAll(".view-content-selector")
  contentSelectors.forEach(selector => {
    selector.style.backgroundColor = "var(--main-background-color"
    selector.style.borderBottom = "solid 1px var(--gpg-green)"
  })

  const clickedSelector = event.currentTarget
  clickedSelector.style.backgroundColor = "white"
  clickedSelector.style.borderBottom = "none"
}