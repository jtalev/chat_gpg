
function toggleContentSelector(event) {
	const contentSelectors = document.querySelectorAll(".view-content-selector")
	contentSelectors.forEach(selector => {
		selector.style.backgroundColor = "var(--main-background-color)"
		selector.style.borderBottom = "solid 1px var(--gpg-green)"
	})

	const target = event.currentTarget
	target.style.backgroundColor = "white"
	target.style.borderBottom = "none"
}

async function sendPurchaseOrder(event) {
	event.preventDefault()

	const form = document.querySelector("#purchase-order-form")

	const employeeId = form.querySelector("#employee_id").value
	const storeId = form.querySelector("#store_id").value
	const jobId = parseInt(form.querySelector("#job_id").value)
	const date = form.querySelector("#date").value

	const itemRows = Array.from(form.querySelectorAll('.purchase-order-item-row'))
	  .slice(1) // skip the header row
	  .map(row => {
	    const quantity = parseInt(row.querySelector('input[name="quantity"]')?.value);
	    const item_type_id = row.querySelector('select[name="item_type_id"]')?.value;
	    const item_size_id = row.querySelector('select[name="item_size_id"]')?.value;
	    const item_name = row.querySelector('input[name="item_name"]')?.value;
	    return { quantity, item_type_id, item_size_id, item_name };
    });
    console.log(itemRows)

	const payload = {
		employee_id: employeeId,
		store_id: storeId,
		job_id: jobId,
		date: date,
		purchase_order_items: itemRows
	}

	try {
      const response = await fetch('/purchase-order/post', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload)
      });

      const result = await response.text();
      const container = document.getElementById('purchase-order-view-content')
      container.innerHTML = result;
      htmx.process(container)
    } catch (error) {
      console.error('Error submitting purchase order:', error);
    }
}

function deleteItemRow(event) {
    const row = event.target.closest('.purchase-order-item-row');
    
    if (row) {
        row.remove();
    }
}

function toggleModal() {
	const modal = document.querySelector(".modal")
	modal.style.display = modal.style.display === "flex" ? "none" : "flex"
}

function toggleSubmitBtn(event) {
	const submitBtn = event.target
	const btnText = submitBtn.textContent

	if (btnText === "SUBMIT") {
		submitBtn.textContent = "CONFIRM ORDER"
		submitBtn.style.backgroundColor = "green"
		submitBtn.style.width = "auto"
	} else if (btnText === "CONFIRM ORDER") {
		sendPurchaseOrder(event)
	} else if (btnText === "Cancel") {
		const btn = document.querySelector("#purchase-order-submit-btn")
		btn.textContent = "SUBMIT"
		btn.style.backgroundColor = "var(--gpg-green)"
		btn.style.width = "80px"
	}
}