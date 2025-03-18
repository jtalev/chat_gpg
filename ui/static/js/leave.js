document.addEventListener("DOMContentLoaded", onLeaveSelectorClick)
function onLeaveSelectorClick() {
    document.querySelectorAll(".leaveContentSelector").forEach(selector => {
        selector.addEventListener("click", function () {
            const leaveContentSelectorChildren = document.querySelectorAll(".leaveContentSelector")
            leaveContentSelectorChildren.forEach(element => {
                element.style.backgroundColor = "var(--main-background-color)"
                element.style.borderBottom = ("solid 1px var(--gpg-green)")
            });

            this.style.backgroundColor = "white"
            this.style.borderBottom = "white"

            let leaveContentElement
            switch (this.textContent.trim()) {
                case "REQUEST":
                    leaveContentElement = document.querySelector("#leaveRequestContent")
                    leaveContentElement.style.display = "flex"
                    break
                case "HISTORY":
                    leaveContentElement = document.querySelector("#leaveHistoryContent")
                    leaveContentElement.style.display = "flex"
                    break
            }
        })
    })
}

document.getElementById("leave-form").addEventListener("htmx:afterSwap", function(event) {
    console.log("after swap")
    onLeaveFormSubmit()
})

function onLeaveFormSubmit() {
    const submitBtn = document.getElementById("leave-form-submit-btn")
    const cancelBtn = document.getElementById("leave-form-cancel-btn")
    const error = document.querySelector(".error")
    let isFormValid = true
    
    if (error.textContent != "") {
        isFormValid = false
    }

    console.log(error.textContent)
    console.log(isFormValid)
    
    if (!isFormValid) {
        return
    }

    const leaveSubmitAlert = document.getElementById("leave-submit-alert-container").style.display = "flex"
    cancelBtn.style.display = "none"
    submitBtn.textContent = "SUBMITTED"
    submitBtn.disabled = true   
}

function onLeaveFormReset() {
    const submitBtn = document.getElementById("leave-form-submit-btn")
    const cancelBtn = document.getElementById("leave-form-cancel-btn")
    const errors = document.querySelector(".error")

    const leaveSubmitAlert = document.getElementById("leave-submit-alert-container").style.display = "none"
    cancelBtn.style.display = "flex"
    submitBtn.textContent = "SUBMIT"
    submitBtn.disabled = false
    errors.textContent = ""
}

document.addEventListener("DOMContentLoaded", function () {
    const dateInputs = document.querySelectorAll(".dateInput")
    dateInputs.forEach(element => {
        element.addEventListener("blur", function () {
            
        })
    })
})


document.addEventListener("htmx:afterSwap", onLeaveHistoryRowClick)
function onLeaveHistoryRowClick() {
    const leaveHistoryRows = document.querySelectorAll(".leaveHistoryRow")
    leaveHistoryRows.forEach(element => {
        element.addEventListener("click", function () {
            if (element.classList.contains("leaveHistoryRowHeader")) return

            const allRowsDisplayed = document.querySelectorAll(".leaveHistoryRowDisplayed")
            const allRowsHidden = document.querySelectorAll(".leaveHistoryRowHidden")
            allRowsDisplayed.forEach(element => {
                element.style.display = "flex"
            })
            allRowsHidden.forEach(element => {
                element.classList.add("hidden")
            })
            
            const displayedRow = element.querySelector(".leaveHistoryRowDisplayed")
            const hiddenRow = element.querySelector(".leaveHistoryRowHidden")
            displayedRow.style.display = "none"
            hiddenRow.classList.remove("hidden")
        })
    })
}

document.addEventListener("htmx:afterSwap", onLeaveHeaderClick);
function onLeaveHeaderClick() {
    const leaveHeaders = document.querySelectorAll(".employee-leave-container-header");
    console.log(leaveHeaders)

    leaveHeaders.forEach(header => {
        header.addEventListener("click", function () {
            console.log("click")
            const leaveRequestContainers = document.querySelectorAll(".employee-leave-list-container");
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
    document.querySelector("#employee-leave-modal-container").style.display = "flex"
}

function onModalClose() {
    document.querySelector("#employee-leave-modal-container").style.display = "none"
}