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

function toggleEmployeeBreakdown(event) {
  const clickedRow = event.currentTarget
  const employeeId = clickedRow.id
  let divToExpand

  const expandables = document.querySelectorAll(".job-report-employee-breakdown")
  expandables.forEach(div => {
    if (div.id === employeeId) {
      divToExpand = div
    }
  })

  if (!divToExpand) return

  divToExpand.style.display = divToExpand.style.display === "flex" ? "none" : "flex"

  rotateArror(clickedRow)
}

function rotateArror(clickedRow) {
  clickedRow.style.rotate = clickedRow.style.rotate === "180deg" ? "0deg" : "180deg"
}