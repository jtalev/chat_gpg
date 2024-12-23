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
                case "TIMESHEET":
                    leaveContentElement = document.querySelector("#adminTimesheetContent")
                    leaveContentElement.style.display = "block"
                    break
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