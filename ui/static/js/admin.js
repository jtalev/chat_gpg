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