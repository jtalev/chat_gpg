document.addEventListener("DOMContentLoaded", onLeaveSelectorClick)
function onLeaveSelectorClick() {
    document.querySelectorAll(".adminContentSelector").forEach(selector => {
        selector.addEventListener("click", function () {

            const leaveContentSelectorChildren = document.querySelectorAll(".adminContentSelector")
            leaveContentSelectorChildren.forEach(element => {
                element.style.backgroundColor = "var(--main-background-color)"
                element.style.borderBottom = ("solid 1px var(--gpg-green)")
            });

            const leaveContentContainerChildren = document.querySelectorAll(".adminContent")
            leaveContentContainerChildren.forEach(element => {
                element.style.display = "none"
            })

            this.style.backgroundColor = "white"
            this.style.borderBottom = "white"

            let leaveContentElement
            switch (this.textContent.trim()) {
                case "JOBS":
                    leaveContentElement = document.querySelector("#adminJobsContent")
                    leaveContentElement.style.display = "block"
                    break
                case "LEAVE":
                    leaveContentElement = document.querySelector("#adminLeaveContent")
                    leaveContentElement.style.display = "block"
                    break
                case "EMPLOYEES":
                    leaveContentElement = document.querySelector("#adminEmployeesContent")
                    leaveContentElement.style.display = "block"
                    break
                case "SAFETY":
                    leaveContentElement = document.querySelector("#adminSafetyContent")
                    leaveContentElement.style.display = "block"
                    break
                case "PURCHASE ORDER":
                    leaveContentElement = document.querySelector("#adminPurchaseOrderContent")
                    leaveContentElement.style.display = "block"
                    break
            }
        })
    })
}

document.addEventListener("htmx:afterSwap", onAddJobClick)
function onAddJobClick() {
    const addJobBtn = document.querySelector("#add-job-btn")
    addJobBtn.addEventListener("click", function () {
        const addJobModalContainer = document.querySelector("#add-job-modal-container")
        addJobModalContainer.style.display = "flex"
    })
}

function onAddJobClose() {
    document.querySelector("#add-job-modal-container").style.display = "none"
}

document.addEventListener("htmx:afterSwap", onAddJobSubmit)
function onAddJobSubmit() {
    const addJobSubmitBtn = document.querySelector("#add-form-submit-btn")
    addJobSubmitBtn.addEventListener("click", function () {
        const addJobModalContainer = document.querySelector("#add-job-modal-container")
        addJobModalContainer.style.display = "none"
    })
}

document.addEventListener("htmx:afterSwap", onLeaveHeaderClick)
function onLeaveHeaderClick() {
    const leaveHeaders = document.querySelectorAll(".admin-leave-headers");

    leaveHeaders.forEach(header => {
        header.addEventListener("click", function () {
            const leaveRequestContainers = document.querySelectorAll(".admin-leave-request-container");
            document.querySelector("#pending-arrow").style.rotate = "-90deg"
            document.querySelector("#approved-arrow").style.rotate = "-90deg"
            document.querySelector("#denied-arrow").style.rotate = "-90deg"
            leaveRequestContainers.forEach(container => {
                container.style.display = "none";
                if (this.id === "pending" && container.id === "leave-pending-requests") {
                    container.style.display = "flex"
                    document.querySelector("#pending-arrow").style.rotate = "0deg"
                }
                if (this.id === "approved" && container.id === "leave-approved-requests") {
                    container.style.display = "flex"
                    document.querySelector("#approved-arrow").style.rotate = "0deg"
                }
                if (this.id === "denied" && container.id === "leave-denied-requests") {
                    container.style.display = "flex"
                    document.querySelector("#denied-arrow").style.rotate = "0deg"
                }
            });
        });
    });
}

function onLeaveRequestClick() {
    document.querySelector("#admin-leave-modal-container").style.display = "flex"
}

function onModalClose() {
    document.querySelector("#admin-leave-modal-container").style.display = "none"
}

function onAddEmployeeClick() {
    document.querySelector("#add-employee-modal-container").style.display = "flex"
}

function onAddEmployeeSubmit() {
    const err = document.querySelector("#add-employee-err")
    if (err.textContent === "") {
        onAddEmployeeClose()
    }
}

function onAddEmployeeClose() {
    document.querySelector("#add-employee-modal-container").style.display = "none"
}

document.addEventListener("DOMContentLoaded", onAdminEmployeeRowClick)
document.addEventListener("htmx:afterSwap", onAdminEmployeeRowClick)
function onAdminEmployeeRowClick() {
    document.querySelectorAll(".admin-employee-row p").forEach(function (p) {
        p.addEventListener("click", onPutEmployeeClick);
    })
}

function onPutEmployeeClick() {
    document.querySelector("#admin-employee-put-modal-container").style.display = "flex"
}

function onPutEmployeeClose() {
    document.querySelector("#admin-employee-put-modal-container").style.display = "none"
}

function onPutEmployeeSuccessful() {
    document.querySelector("#admin-employee-put-modal-container").style.display = "none"
}

function onPutJobClick() {
    document.querySelector("#put-job-modal-container").style.display = "flex"
}

function onPutJobModalClose() {
    document.querySelector("#put-job-modal-container").style.display = "none"
}

function onSafetyContentToggle(event) {
    // Get toggle elements
    const incidentToggle = document.getElementById("incident-report-toggle");
    const swmToggle = document.getElementById("swm-toggle");

    // Reset both toggles to default
    incidentToggle.style.backgroundColor = "var(--gpg-green-background)";
    incidentToggle.style.color = "white";
    swmToggle.style.backgroundColor = "var(--main-background-color)";
    swmToggle.style.color = "black";

    // Highlight the clicked toggle
    if (event.target.id === "incident-report-toggle") {
        incidentToggle.style.backgroundColor = "var(--gpg-green-background)";
        incidentToggle.style.color = "white";
        swmToggle.style.backgroundColor = "var(--main-background-color)";
        swmToggle.style.color = "black";
    } else if (event.target.id === "swm-toggle") {
        swmToggle.style.backgroundColor = "var(--gpg-green-background)";
        swmToggle.style.color = "white";
        incidentToggle.style.backgroundColor = "var(--main-background-color)";
        incidentToggle.style.color = "black";
    }
}

function togglePurchaseOrderContentSelector(event) {
    const target = event.currentTarget
    const selectors = document.querySelectorAll(".admin-purchase-order-content-selector")
    selectors.forEach(selector => {
        selector.style.backgroundColor = "var(--light-background-color)"
        selector.style.color = "var(--timesheet-border)"
    })

    target.style.backgroundColor = "var(--gpg-light-green-background)"
    target.style.color = "var(--main-color)"
}

function toggleAdminPurchaseOrderModal() {
    const container = document.getElementById("admin-po-modal-container")
    container.style.display = container.style.display === "flex" ? "none" : "flex"
}

document.body.addEventListener("htmx:afterSwap", function(event) {
    const srcElementId = event.srcElement.id
    if (srcElementId === "add-store-form") {
        htmx.ajax("GET", "/admin/serve-stores-template", {
            target: "#admin-purchase-order-content",
            swap: "innerHTML"
        });
    } else if (srcElementId === "add-item-form") {
        htmx.ajax("GET", "/admin/serve-item-types-template", {
            target: "#admin-purchase-order-content",
            swap: "innerHTML"
        });
    }
})

function toggleModal() {
    const modal = document.querySelector("#admin-po-modal-container")
    modal.style.display = modal.style.display === "flex" ? "none" : "flex"
}

document.addEventListener("htmx:afterSwap", function(event) {
    const modalId = event.srcElement.children[0].id

    if (modalId === "item-modal") {
        fetch('/admin/serve-item-types-template')
        .then(response => {
          if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
          }
          return response.text();
        })
        .then(html => {
            const container = document.querySelector('#admin-purchase-order-content')
            container.innerHTML = html;
        })
        .catch(error => {
          console.error('Error fetching item types template:', error);
        });
    } else if (modalId === "store-modal") {
        fetch('/admin/serve-stores-template')
        .then(response => {
          if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
          }
          return response.text();
        })
        .then(html => {
            const container = document.querySelector('#admin-purchase-order-content')
            container.innerHTML = html;
        })
        .catch(error => {
          console.error('Error fetching item types template:', error);
        });
    }

})