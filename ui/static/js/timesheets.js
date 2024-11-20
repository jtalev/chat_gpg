function getDateContainerChildren() {
    const dateContainer = document.getElementById("date-container")
    const dateContainerChildren = dateContainer.children
    return dateContainerChildren
}

function getPreviousWednesday(currentDate) {
    const dayOfWeek = currentDate.getDay()
    const daysSinceWednesday = (dayOfWeek + 4) % 7
    const dateToFill = currentDate.getDate() - daysSinceWednesday

    const previousWednesday = new Date(currentDate)
    previousWednesday.setDate(dateToFill)

    return previousWednesday
}

function fillDateCell() {
    const currentDate = new Date()
    const dateContainerChildren = getDateContainerChildren()
    let previousWednesday = getPreviousWednesday(currentDate)

    for (var i = 0; i < dateContainerChildren.length; i++) {
        const child = dateContainerChildren[i]
        const dateElement = child.querySelector(".date")
        dateElement.textContent = previousWednesday.getDate()
        if (previousWednesday.getDate() === currentDate.getDate()) {
            dateElement.style.backgroundColor = "#3079d1"
            dateElement.style.color = "white"
        }
        previousWednesday.setDate(previousWednesday.getDate() + 1)
    }
}

document.querySelectorAll(".date-cell").forEach(cell => {
    cell.addEventListener("click", function () {
        const dateContainerChildren = getDateContainerChildren()

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

document.querySelector("#left-arrow-container").addEventListener("click", function () {
    const wedDateElement = document.getElementById("wed-date")
    const payweekMonthElement = document.getElementById("payweek-month")
    const payweekYearElement = document.getElementById("payweek-year")

    monthNames = [
        "January", "February", "March", "April", "May", "June",
        "July", "August", "September", "October", "November", "December"
    ]
    let payweekMonth
    for (i = 0; i < monthNames.length; i++) {
        if (monthNames[i] === payweekMonthElement.innerHTML) {
            payweekMonth = i
        }
    }
    
    const wedDate = Number(wedDateElement.innerHTML)
    const payweekYear = Number(payweekYearElement.innerHTML)
    
    let currentWed = new Date()
    currentWed.setFullYear(payweekYear)
    currentWed.setMonth(payweekMonth)
    currentWed.setDate(wedDate - 7)

    if (currentWed.getMonth() < payweekMonth) {
        payweekMonthElement.innerHTML = monthNames[currentWed.getMonth()]
    }

    const dateContainerChildren = getDateContainerChildren()

    for (var i = 0; i < dateContainerChildren.length; i++) {
        const child = dateContainerChildren[i]
        const dateElement = child.querySelector(".date")
        dateElement.textContent = currentWed.getDate()
        currentWed.setDate(currentWed.getDate() + 1)
    }
})

function fillMonth() {
    const payweekMonthElement = document.getElementById("payweek-month")
    const payweekYearElement = document.getElementById("payweek-year")
    
    const date = new Date()
    const currentYear = date.getFullYear()
    const currentMonth = date.getMonth()
    monthNames = [
        "January", "February", "March", "April", "May", "June",
        "July", "August", "September", "October", "November", "December"
    ]

    const currentMonthName = monthNames[currentMonth]

    payweekMonthElement.innerHTML = currentMonthName
    payweekYearElement.innerHTML = currentYear
}

document.addEventListener("DOMContentLoaded", fillDateCell)
document.addEventListener("DOMContentLoaded", fillMonth)