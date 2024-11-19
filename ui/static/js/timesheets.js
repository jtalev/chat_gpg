function fillDateCell() {
    const dateContainer = document.getElementById("date-container")
    const dateContainerChildren = dateContainer.children

    const date = new Date()

    const dayOfWeek = date.getDay()
    const daysSinceWednesday = (dayOfWeek + 4) % 7
    let dateToFill = date.getDate() - daysSinceWednesday

    for (var i = 0; i < dateContainerChildren.length; i++) {
        const child = dateContainerChildren[i]
        const dateElement = child.querySelector(".date")
        dateElement.textContent = dateToFill
        if (dateToFill === date.getDate()) {
            dateElement.style.backgroundColor = "#3079d1"
            dateElement.style.color = "white"
        }
        dateToFill++
    }
}

document.querySelectorAll(".date-cell").forEach(cell => {
    cell.addEventListener("click", function () {
        const dateContainer = document.getElementById("date-container")
        const dateContainerChildren = dateContainer.children

        for (var i = 0; i < dateContainerChildren.length; i++) {
            const child = dateContainerChildren[i]
            let dateElement = child.querySelector(".date")
            dateElement.style.backgroundColor = "#cccccc"
            dateElement.style.color = "black"
        }

        let dateElement = this.querySelector(".date")
        dateElement.style.backgroundColor = "#3079d1"
        dateElement.style.color = "white"
    })
})

document.addEventListener("DOMContentLoaded", fillDateCell)