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
