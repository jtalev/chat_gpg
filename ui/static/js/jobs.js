function toggleNoteModal() {
	const modal = document.getElementById("note-modal")
	const display = getComputedStyle(modal).display;
	modal.style.display = display == "none" ? "flex" : "none";
}

async function submitPaintnote(event) {
	event.preventDefault()

	const form = document.getElementById("paintnote-form")
	const formData = new FormData(form)
	const jsonData = {}

	formData.forEach((value, key) => {
		jsonData[key] = value
	})

	if (jsonData.job_id) {
		jsonData.job_id = parseInt(jsonData.job_id)
	}
	if (jsonData.coats) {
		jsonData.coats = parseInt(jsonData.coats)
	} else {
		jsonData.coats = 0
	}

	const formType = form.getAttribute("data-form-type")
	const endpoint = formType === "post"
		? "/job-notes/post"
		: "/job-notes/put"
	const method = formType === "post"
		? "POST"
		: "PUT"

	try {
		const response = await fetch(endpoint, {
			method: method,
			headers: {
				"Content-Type": "application/json",
			},
			body: JSON.stringify(jsonData)
		})

		if (!response.ok) {
			const errorText = await response.text()
			console.error("server error:", errorText)
			alert("failed to submit paint note")
			return false
		}

		const html = await response.text()
		document.getElementById("note-modal").innerHTML = html
		return false
	} catch (err) {
		console.log(err)
		alert("error occurred while submitting paint note")
		return false
	}
}

async function submitImagenote(event) {
	event.preventDefault()

	const form = document.getElementById("imagenote-form")
	const formData = new FormData(form)
	const jsonData = {}

	formData.forEach((value, key) => {
		jsonData[key] = value
	})

	if (jsonData.job_id) {
		jsonData.job_id = parseInt(jsonData.job_id)
	}

	const fileInput = document.querySelector('input[name="image_base64"]')
	const file = fileInput.files[0]

	if (file) {
		const reader = new FileReader()

		const base64 = await new Promise((resolve, reject) => {
			reader.onload = () => resolve(reader.result.split(",")[1])
			reader.onerror = reject
			reader.readAsDataURL(file)
		})

		jsonData.image_base64 = base64
	}

	console.log(jsonData.image_base64)

	const formType = form.getAttribute("data-form-type")
	const endpoint = formType === "post" 
		? "/job-notes/post"
		: "/job-notes/put"
	const method = formType === "post"
		? "POST"
		: "PUT"

	try {
		const response = await fetch(endpoint, {
			method: method,
			headers: {
				"Content-Type": "application/json",
			},
			body: JSON.stringify(jsonData)
		})

		if (!response.ok) {
			const errorText = await response.text()
			console.error("server error:", errorText)
			alert("Failed to submit image note")
			return false
		}

		const html = await response.text()
		document.getElementById("note-modal").innerHTML = html
		return false
	} catch (err) {
		console.error(err)
		alert("Error occurred while submitting image note")
		return false
	}
}
