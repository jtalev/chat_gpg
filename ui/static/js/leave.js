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