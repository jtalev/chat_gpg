document.addEventListener("DOMContentLoaded", onLeaveSelectorClick)
function onLeaveSelectorClick() {
    document.querySelectorAll(".leaveContentSelector").forEach(selector => {
        selector.addEventListener("click", function () {
            const leaveContentSelectorChildren = document.querySelectorAll(".leaveContentSelector")
            leaveContentSelectorChildren.forEach(element => {
                element.style.backgroundColor = "var(--main-background-color)"
                element.style.borderBottom = ("solid 1px var(--gpg-green)")
            });

            const leaveContentContainerChildren = document.querySelectorAll(".leaveContent")
            leaveContentContainerChildren.forEach(element => {
                element.style.display = "none"
            })

            this.style.backgroundColor = "white"
            this.style.borderBottom = "white"

            let leaveContentElement
            switch (this.textContent.trim()) {
                case "REQUEST":
                    leaveContentElement = document.querySelector("#leaveRequestContent")
                    leaveContentElement.style.display = "block"
                    break
                case "HISTORY":
                    leaveContentElement = document.querySelector("#leaveHistoryContent")
                    leaveContentElement.style.display = "block"
                    break
            }
        })
    })
}

document.addEventListener("htmx:afterSwap", onLeaveFormSubmit)
function onLeaveFormSubmit() {
    const submitBtn = document.getElementById("leaveFormSubmitBtn")
    const cancelBtn = document.getElementById("leaveFormCancelBtn")
    const errors = document.querySelectorAll(".error")
    let isFormValid = true
    errors.forEach(error => {
        if (error.textContent != "") {
            isFormValid = false
        }
    })
    if (!isFormValid) {
        return
    }

    cancelBtn.textContent = "Reset"
    submitBtn.textContent = "SUBMITTED"
    submitBtn.disabled = true   
}

document.addEventListener("DOMContentLoaded", onLeaveFormReset)
function onLeaveFormReset() {
    const cancelBtn = document.getElementById("leaveFormCancelBtn")
    cancelBtn.addEventListener("click", function () {
        const submitBtn = document.getElementById("leaveFormSubmitBtn")
    
        cancelBtn.textContent = "Cancel"
        submitBtn.textContent = "SUBMIT"
        submitBtn.disabled = false
    })
}

document.addEventListener("DOMContentLoaded", function () {
    const dateInputs = document.querySelectorAll(".dateInput")
    dateInputs.forEach(element => {
        element.addEventListener("blur", function () {
            
        })
    })
})

document.addEventListener("DOMContentLoaded", function () {
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
})