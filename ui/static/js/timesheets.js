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
    const tableBody = document.getElementById("timesheet-table-job-rows")
    const table = document.getElementById("timesheetTable")
    if (!table || !tableBody) return

    const rows = tableBody.querySelectorAll("tr")
    const numRows = rows.length
    const dayTotalRow = document.querySelectorAll(".day-total")

    if (numRows === 0) {
        console.log("No data rows, resetting totals")
        dayTotalRow.forEach(cell => {
            cell.innerHTML = "0:0"
        })
        return
    }

    const numCols = table.querySelectorAll(".timeInput").length / numRows

    for (let colIdx = 0; colIdx < numCols; colIdx++) {
        let totalHours = 0
        let totalMinutes = 0

        rows.forEach(row => {
            const cell = row.querySelectorAll(".timeInput")[colIdx]
            if (cell) {
                let cellValue = cell.value.trim()
                if (cellValue === "") {
                    cellValue = "0:0"
                }

                if (cellValue.includes(":")) {
                    const [hrs, mins] = cellValue.split(":").map(Number)
                    totalHours += isNaN(hrs) ? 0 : hrs
                    totalMinutes += isNaN(mins) ? 0 : mins
                } else {
                    totalHours += isNaN(Number(cellValue)) ? 0 : Number(cellValue)
                }
            }
        })

        totalHours += Math.floor(totalMinutes / 60)
        totalMinutes = totalMinutes % 60

        if (dayTotalRow[colIdx]) {
            dayTotalRow[colIdx].innerHTML = `<p>${totalHours}:${totalMinutes.toString().padStart(1, "0")}</p>`
        }
    }
}

function calculateWeekTotal() {
    const dayTotalCells = document.querySelectorAll(".day-total")
    const rowTotalCells = document.querySelectorAll(".total")
    const tableTotalCell = document.getElementById("timesheet-table-total")

    let rowHrs = 0
    let rowMins = 0
    let dayHrs = 0
    let dayMins = 0

    if (rowTotalCells) {
        rowTotalCells.forEach(cell => {
            cellValue = cell.textContent.trim()
            if (cellValue.includes(":")) {
                const [hrs, mins] = cellValue.split(":").map(Number)
                rowHrs += isNaN(hrs) ? 0 : hrs
                rowMins += isNaN(mins) ? 0 : mins
            } else {
                rowHrs += isNaN(cellValue) ? 0 : Number(cellValue)
            }
        })
    }

    rowHrs += Math.floor(rowMins / 60)
    rowMins = rowMins % 60

    dayTotalCells.forEach(cell => {
        console.log(cell.textContent)
        cellValue = cell.textContent.trim()
        if (cellValue.includes(":")) {
            const [hrs, mins] = cellValue.split(":").map(Number)
            dayHrs += isNaN(hrs) ? 0 : hrs
            dayMins += isNaN(mins) ? 0 : mins
        } else {
            dayHrs += isNaN(cellValue) ? 0 : Number(cellValue)
        }  
    })

    dayHrs += Math.floor(dayMins / 60)
    dayMins = dayMins % 60

    if (rowTotalCells && dayHrs === rowHrs && dayMins === rowMins) {
        tableTotalCell.innerHTML = `${rowHrs}:${rowMins}`
    } else {
        tableTotalCell.innerHTML = `${dayHrs}:${dayMins}`
    }
}

function updateRowTotals() {
    document.querySelectorAll(".timeInput").forEach(cell => {
        cell.removeEventListener("keyup", calculateRowTotals)
        cell.addEventListener("keyup", calculateRowTotals)
        cell.removeEventListener("keyup", calculateDayTotals)
        cell.addEventListener("keyup", calculateDayTotals)
        cell.removeEventListener("keyup", calculateWeekTotal)
        cell.addEventListener("keyup", calculateWeekTotal)
    });
}

document.addEventListener("DOMContentLoaded", () => {
    updateRowTotals()
    calculateRowTotals()
    calculateDayTotals()
    calculateWeekTotal()
});

document.body.addEventListener("htmx:afterSwap", function() {
    updateRowTotals()
    calculateRowTotals()
    calculateDayTotals()
    calculateWeekTotal()
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