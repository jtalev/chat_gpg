function togglePdfViewer() {
    if (window.innerWidth > 1029) {
        return;
    }

    const adminSafetyContent = document.getElementById("admin-safety-content");
    const adminSafetyPdfViewer = document.getElementById("admin-safety-pdf-viewer");
    const closePdfBtn = document.getElementById("close-pdf-btn")

    const isContentVisible = window.getComputedStyle(adminSafetyContent).display !== "none";

    if (isContentVisible) {
        adminSafetyContent.style.display = "none";
        adminSafetyPdfViewer.style.display = "flex";
        closePdfBtn.style.display = "block"
    } else {
        adminSafetyContent.style.display = "flex";
        adminSafetyPdfViewer.style.display = "none";
        closePdfBtn.style.display = "none"
    }
}

document.addEventListener("DOMContentLoaded", function() {
    document.querySelectorAll(".incident-report-list-item-btn img[src*='view.svg']").forEach(button => {
        button.addEventListener("click", togglePdfViewer);
    });
});