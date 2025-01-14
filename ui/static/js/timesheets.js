document.addEventListener("htmx:afterSwap", timeInputFilter)
document.addEventListener("DOMContentLoaded", timeInputFilter) 
function timeInputFilter() {
    const timeInputs = document.querySelectorAll('.timeInput')
    
    timeInputs.forEach(timeInput => {
        timeInput.addEventListener('keydown', (event) => {
            // Allow numbers, backspace, delete, colon, and number pad keys
            if (!((event.key >= '0' && event.key <= '9') ||
            event.key === 'Backspace' ||
            event.key === 'Delete' ||
            event.key === ':' ||
            (event.key >= 'NumPad0' && event.key <= 'NumPad9'))) {
                event.preventDefault();
            }
        })
    })
}

document.addEventListener("htmx:afterSwap", calculateRowTotals)
document.addEventListener("DOMContentLoaded", calculateRowTotals)
function calculateRowTotals() {
    const table = document.getElementById("timesheetTable")
    const rows = table.rows

    for (let i = 1; i < rows.length; i++) {
        const cells = rows[i].querySelectorAll(".timeInput")
        let hours = 0
        let mins = 0

        for (let j = 0; j < cells.length; j++) {
            const cellValue = cells[j].value
            const containsColon = cellValue.includes(":")
            if (containsColon) {
                const cellValueSplit = cellValue.split(":")
                const addHrs = cellValueSplit[0]
                const addMins = cellValueSplit[1]
                hours += Number(addHrs)
                mins += Number(addMins)
            } else {
                hours += Number(cellValue)
            }
        }

        // convert mins to hours and minutes and add to totals
        const extraHours = Math.floor(mins / 60)
        hours += extraHours
        mins = mins % 60

        const total = rows[i].querySelector(".total")
        total.innerHTML = hours + ":" + mins
    }
}

document.addEventListener("DOMContentLoaded", updateRowTotals)
document.addEventListener("keyup", updateRowTotals)
function updateRowTotals() {
    document.querySelectorAll(".timeInput").forEach(cell => {
        cell.removeEventListener("keyup", calculateRowTotals)
        cell.addEventListener("keyup", calculateRowTotals)
    })
}

document.addEventListener("DOMContentLoaded", onAddRowClick)
document.addEventListener("htmx:afterSwap", onAddRowClick)
function onAddRowClick() {
    document.querySelector("#addTableRowBtn").addEventListener("click", function() {
        const jobSelectModal = document.querySelector("#jobSelectModal")
        jobSelectModal.style.display = "flex"
    })
}

document.addEventListener("DOMContentLoaded", onJobSelectSubmit)
document.addEventListener("htmx:afterSwap", onJobSelectSubmit)
function onJobSelectSubmit() {
    document.querySelector("#selectJobSubmitBtn").addEventListener("click", function() {
        const jobSelectModal = document.querySelector("#jobSelectModal")
        jobSelectModal.style.display = "none"
    })
}

document.addEventListener("DOMContentLoaded", onJobSelectClose)
document.addEventListener("htmx:afterSwap", onJobSelectClose)
function onJobSelectClose() {
    document.querySelector("#jobSelectModalCloseBtn").addEventListener("click", function() {
        const jobSelectModal = document.querySelector("#jobSelectModal")
        jobSelectModal.style.display = "none"
    })
}