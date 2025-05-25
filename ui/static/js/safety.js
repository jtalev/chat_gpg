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

function toggleContentSelector(event) {
    const clickedSelector = event.currentTarget
    const selectors = document.querySelectorAll(".view-content-selector")

    selectors.forEach(selector => {
        selector.style.backgroundColor = "var(--main-background-color)"
        selector.style.borderBottom = "solid 1px var(--gpg-green)"
        if (selector === clickedSelector) {
            selector.style.backgroundColor = "white"
            selector.style.borderBottom = "none"
        }
    })
}

function toggleSwmViewer() {
    if (window.innerWidth > 1030) {
        return 
    }

    const swmList = document.querySelector("#swm-list-container")
    const swmViewer = document.querySelector("#swm-viewer-container")
    const closePdfBtn = document.getElementById("swms-close-pdf-btn")

    const isContentVisible = window.getComputedStyle(swmList).display !== "none";

    if (isContentVisible) {
        swmList.style.display = "none"
        swmViewer.style.display = "flex"
        closePdfBtn.style.display = "block"
    } else {
        swmList.style.display = "flex"
        swmViewer.style.display = "none"
        closePdfBtn.style.display = "none"
    }
}

document.body.addEventListener("htmx:afterSwap", function(evt) {
    const target = evt.detail.target;
    
    if (target.id === "swm-viewer-container") {
        htmx.ajax("GET", "/safety/swms/get-list-html", {
            target: "#swm-list-container",
            swap: "innerHTML"
        });
    }
});