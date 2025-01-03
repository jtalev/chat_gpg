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


function getDateContainerChildren() {
    const dateContainer = document.getElementById("datesRow")
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

// document.addEventListener("DOMContentLoaded", renderDateCell)
// function renderDateCell() {
//     const currentDate = new Date()
//     const dateContainerChildren = getDateContainerChildren()
//     let previousWednesday = getPreviousWednesday(currentDate)
//     let columnIndex = 0
    
//     // i < datecontainerChildren.length - 1 here because the 
//     // table has an extra column for total hours worked
//     for (var i = 1; i < dateContainerChildren.length - 1; i++) {
//         const child = dateContainerChildren[i]
//         const dateElement = child.querySelector(".date")
//         dateElement.textContent = previousWednesday.getDate()
//         if (previousWednesday.getDate() === currentDate.getDate()) {
//             dateElement.style.color = "white"
//             columnIndex = i
//         }
//         previousWednesday.setDate(previousWednesday.getDate() + 1)
//     }
    
//     const timesheetTable = document.querySelector("#timesheetTable")
//     const timesheetRows = timesheetTable.rows
    
//     // for (let i = 0; i < timesheetRows.length; i++) {
//     //     const row = timesheetRows[i]
//     //     const cells = row.cells
        
//     //     cells[columnIndex].style.backgroundColor = 'gray'
//     //     cells[columnIndex].style.pointerEvents = "all"
//     // }
// }

document.addEventListener("DOMContentLoaded", renderMonth)
function renderMonth() {
    const payweekMonthElement = document.getElementById("payweekMonth")
    const payweekYearElement = document.getElementById("payweekYear")
    
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

document.addEventListener("htmx:afterSwap", addTableRow)
document.addEventListener("DOMContentLoaded", addTableRow)
function addTableRow() {
    document.querySelector("#addTableRowBtn").addEventListener("click", function() {
        const table = document.querySelector("#timesheetTable")
        const newRow = document.createElement('tr')
        newRow.classList.add("projectRow")
    
        const oldRow = document.querySelector("#newProjectRow")
        const rowString = oldRow.innerHTML
        
        newRow.innerHTML = rowString
        newRow.querySelector(".total").innerHTML = "0:0"
        table.appendChild(newRow)

        updateRowTotals()
        timeInputFilter()
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

function putTimesheets() {
    const forms = document.querySelectorAll(".saveAllForm")
    const timesheetData = []

    forms.forEach((form) => {
        const timeInput = form.querySelectorAll("input[name='time']")
        timeInput.forEach(input => {
            if (input.id != 0) {
                timesheetData.push({
                    id: input.id,
                    time: input.value
                })
            }
        })
    })

    fetch('/timesheets/put-all', {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ timesheets: timesheetData })
    }).then((response => {
        if (response.ok) {
            console.log('Timesheets updated successfully')
        } else {
            console.error('Failed to update timesheets')
        }
    })).catch((error) => {
        console.error('Error:', error)
    })
}

function postTimesheets() {
    const forms = document.querySelectorAll(".saveAllForm");
    const timesheetData = [];
    const weekStartDate = document.querySelector("#wedDate").innerText;
    const month = document.querySelector("#payweekMonth").innerText;
    const year = document.querySelector("#payweekYear").innerText;
    const tableRows = document.getElementById("timesheetTable").rows;

    for (let i = 1; i < tableRows.length; i++) {
        const cells = tableRows[i].cells;
        for (let j = 1; j < cells.length - 2; j++) {
            const timeInput = cells[j].querySelectorAll("input[name='time']");
            timeInput.forEach(input => {
                if (input.id == 0 && input.value !== "") {
                    const job = cells[0].querySelector(".projects");
                    let jobId = job.dataset.jobid;
                    if (jobId == undefined) {
                        const selectedOption = job.options[job.selectedIndex]
                        jobId = selectedOption.getAttribute('data-jobId')
                        console.log("getting selected option")
                    }
                    const tsDate = tableRows[0].cells[j].querySelector(".date").innerText
                    timesheetData.push({
                        job: jobId,
                        time: input.value,
                        weekStart: {
                            date: weekStartDate,
                            month: month,
                            year: year
                        },
                        date: tsDate
                    });
                }
            });
        }
    }
    console.log(timesheetData)
    
    fetch('/timesheets/post-all', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ timesheets: timesheetData })
    }).then((response => {
        if (response.ok) {
            console.log('Timesheets inserted successfully')
        } else {
            console.error('Failed to post timesheets')
        }
    })).catch((error) => {
        console.error('Error:', error)
    })
}


document.addEventListener("htmx:afterSwap", saveEventListenerMount)
document.addEventListener("DOMContentLoaded", saveEventListenerMount) 
function saveEventListenerMount() {
    document.getElementById('saveAllBtn').addEventListener("click", function() {
        putTimesheets()
        postTimesheets()
    })
}