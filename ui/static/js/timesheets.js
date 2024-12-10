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

document.addEventListener("DOMContentLoaded", renderDateCell)
function renderDateCell() {
    const currentDate = new Date()
    const dateContainerChildren = getDateContainerChildren()
    let previousWednesday = getPreviousWednesday(currentDate)
    let columnIndex = 0
    
    // i < datecontainerChildren.length - 1 here because the 
    // table has an extra column for total hours worked
    for (var i = 1; i < dateContainerChildren.length - 1; i++) {
        const child = dateContainerChildren[i]
        const dateElement = child.querySelector(".date")
        dateElement.textContent = previousWednesday.getDate()
        if (previousWednesday.getDate() === currentDate.getDate()) {
            dateElement.style.backgroundColor = "var(--gpg-green)"
            dateElement.style.color = "white"
            columnIndex = i
        }
        previousWednesday.setDate(previousWednesday.getDate() + 1)
    }
    
    const timesheetTable = document.querySelector("#timesheetTable")
    const timesheetRows = timesheetTable.rows
    
    // for (let i = 0; i < timesheetRows.length; i++) {
    //     const row = timesheetRows[i]
    //     const cells = row.cells
        
    //     cells[columnIndex].style.backgroundColor = 'gray'
    //     cells[columnIndex].style.pointerEvents = "all"
    // }
}

// document.addEventListener("DOMContentLoaded", onDateCellClick)
// function onDateCellClick() {
//     document.querySelectorAll(".dateCell").forEach(cell => {
//         cell.addEventListener("click", function () {
//             const dateContainerChildren = getDateContainerChildren()
//             let columnIndex = 0
            
//             // i < datecontainerChildren.length - 1 here because the 
//             // table has an extra column for total hours worked
//             for (var i = 1; i < dateContainerChildren.length - 1; i++) {
//                 const child = dateContainerChildren[i]
//                 let dateElement = child.querySelector(".date")
//                 dateElement.style.backgroundColor = "#cccccc"
//                 dateElement.style.color = "var(--main-color)"
                
//                 let dateCell = child.querySelector(".dateCell")
//                 if (dateCell === this) {
//                     columnIndex = i
//                 }
//             }
            
//             let dateElement = this.querySelector(".date")
//             dateElement.style.backgroundColor = "var(--gpg-green)"
//             dateElement.style.color = "white"
            
//             // colour background of column
            
//             const timesheetTable = document.querySelector("#timesheetTable")
//             const timesheetRows = timesheetTable.rows
            
//             for (let i = 0; i < timesheetRows.length; i++) {
//                 const row = timesheetRows[i]
//                 const cells = row.cells
                
//                 for (let i = 0; i < cells.length; i++) {
//                     cells[i].style.backgroundColor = 'white'
//                 }

//                 if (i > 0) {
//                     for (let j = 1; j < cells.length; j++) {
//                         cells[j].style.pointerEvents = "none"
//                     }
//                 }
                
//                 cells[columnIndex].style.backgroundColor = 'gray'
//                 cells[columnIndex].style.pointerEvents = "all"
//             }
//         })
//     })
// }

document.addEventListener("DOMContentLoaded", onLeftArrowClick)
function onLeftArrowClick() {
    // updates the date carousel, starting from the wednesday previous to the current week
    // converts html to Date() type to update the html again with the correct values
    document.querySelector("#leftArrowContainer").addEventListener("click", function () {
        const wedDateElement = document.getElementById("wedDate")
        const payweekMonthElement = document.getElementById("payweekMonth")
        const payweekYearElement = document.getElementById("payweekYear")
        
        // need to convert month name from string to int; currentWed.setMonth(payweekMonth)
        // use monthNames to update html to new month also
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
        
        // update month names and year
        if (currentWed.getMonth() == 11 && payweekMonth == 0) {
            payweekMonthElement.innerHTML = monthNames[currentWed.getMonth()]
        }
        if (currentWed.getMonth() < payweekMonth) {
            payweekMonthElement.innerHTML = monthNames[currentWed.getMonth()]
        }
        payweekYearElement.innerHTML = currentWed.getFullYear()
        
        // update all dates in carousel div
        const dateContainerChildren = getDateContainerChildren()
        
        for (var i = 1; i < dateContainerChildren.length - 1; i++) {
            const child = dateContainerChildren[i]
            const dateElement = child.querySelector(".date")
            dateElement.textContent = currentWed.getDate()
            currentWed.setDate(currentWed.getDate() + 1)
        }
        
    })
}

document.addEventListener("DOMContentLoaded", onRightArrowClick)
function onRightArrowClick() {
    // updates the date carousel, starting from the wednesday previous to the current week
    // converts html to Date() type to update the html again with the correct values
    document.querySelector("#rightArrowContainer").addEventListener("click", function () {
        const wedDateElement = document.getElementById("wedDate")
        const payweekMonthElement = document.getElementById("payweekMonth")
        const payweekYearElement = document.getElementById("payweekYear")

        // need to convert month name from string to int; currentWed.setMonth((int)payweekMonth)
        // use monthNames to update html to new month also
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
        currentWed.setDate(wedDate + 7)

        // update month names and year
        if (currentWed.getMonth() == 0 && payweekMonth == 11) {
            payweekMonthElement.innerHTML = monthNames[currentWed.getMonth()]
        }
        if (currentWed.getMonth() > payweekMonth) {
            payweekMonthElement.innerHTML = monthNames[currentWed.getMonth()]
        }
        payweekYearElement.innerHTML = currentWed.getFullYear()

        // update all dates in carousel div
        const dateContainerChildren = getDateContainerChildren()

        for (var i = 1; i < dateContainerChildren.length - 1; i++) {
            const child = dateContainerChildren[i]
            const dateElement = child.querySelector(".date")
            dateElement.textContent = currentWed.getDate()
            currentWed.setDate(currentWed.getDate() + 1)
        }
    })
}

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

document.addEventListener("DOMContentLoaded", addTableRow)
function addTableRow() {
    document.querySelector("#addTableRowBtn").addEventListener("click", function() {
        const table = document.querySelector("#timesheetTable")
        const newRow = document.createElement('tr')
        newRow.classList.add("projectRow")
    
        const oldRow = document.querySelector("#projectRow")
        const rowString = oldRow.innerHTML
        
        newRow.innerHTML = rowString
        newRow.querySelector(".total").innerHTML = "0:0"
        table.appendChild(newRow)
    })
}

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
function updateRowTotals() {
    document.querySelectorAll(".timeInput").forEach(cell => {
        cell.addEventListener("keyup", function () {
            calculateRowTotals()
        })
    })
}