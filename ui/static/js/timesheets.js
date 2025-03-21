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

function calculateRowTotals() {
    const table = document.getElementById("timesheetTable")
    if (!table) return
    const rows = table.rows

    for (let i = 1; i < rows.length; i++) {
        const cells = rows[i].querySelectorAll(".timeInput")
        let hours = 0
        let mins = 0

        cells.forEach(cell => {
            const cellValue = cell.value
            if (cellValue.includes(":")) {
                const [addHrs, addMins] = cellValue.split(":").map(Number)
                hours += addHrs
                mins += addMins
            } else {
                hours += Number(cellValue)
            }
        });

        hours += Math.floor(mins / 60)
        mins = mins % 60

        const total = rows[i].querySelector(".total")
        if (total) {
            total.innerHTML = `${hours}:${mins}`
        }
    }
}

function calculateDayTotals() {
    const table = document.getElementById("timesheetTable")
    if (!table) return

    const rows = table.rows
    const numRows = rows.length
    const dayTotalRow = rows[numRows - 1].querySelectorAll(".day-total")
    console.log(numRows)
    if (numRows <= 2) {
        console.log("condition met")
        dayTotalRow.forEach(cell => {
            cell.innerHTML = "0:0"
        })
        return
    }

    const numCols = rows[1].querySelectorAll(".timeInput").length

    for (let colIdx = 0; colIdx < numCols; colIdx++) {
        let totalHours = 0
        let totalMinutes = 0

        // Loop through each row (excluding header & total row)
        for (let rowIdx = 1; rowIdx < numRows - 1; rowIdx++) {
            const cell = rows[rowIdx].querySelectorAll(".timeInput")[colIdx]
            if (cell) {
                let cellValue = cell.value.trim()
                if (cellValue === "") {
                    cellValue = "0:00"  // Treat empty cells as 0 hours, 0 minutes
                }

                if (cellValue.includes(":")) {
                    const [hrs, mins] = cellValue.split(":").map(Number)
                    totalHours += isNaN(hrs) ? 0 : hrs
                    totalMinutes += isNaN(mins) ? 0 : mins
                } else {
                    totalHours += isNaN(Number(cellValue)) ? 0 : Number(cellValue)
                }
            }
        }

        totalHours += Math.floor(totalMinutes / 60)
        totalMinutes = totalMinutes % 60

        if (dayTotalRow[colIdx]) {
            dayTotalRow[colIdx].innerHTML = `<p>${totalHours}:${totalMinutes.toString().padStart(2, "0")}</p>`
        }
    }
}


function updateRowTotals() {
    document.querySelectorAll(".timeInput").forEach(cell => {
        cell.removeEventListener("keyup", calculateRowTotals)
        cell.addEventListener("keyup", calculateRowTotals)
        cell.removeEventListener("keyup", calculateDayTotals)
        cell.addEventListener("keyup", calculateDayTotals)
    });
}

document.addEventListener("DOMContentLoaded", () => {
    updateRowTotals()
    calculateRowTotals()
    calculateDayTotals()
});

document.body.addEventListener("htmx:afterSwap", function() {
    updateRowTotals()
    calculateRowTotals()
    calculateDayTotals()
});

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