package handlers

import (
	"net/http"
	"strconv"

	"github.com/jtalev/chat_gpg/models"
	"github.com/jtalev/chat_gpg/repository"
)

type JobsData struct {
}

func getJobsData() []JobsData {
	data := []JobsData{
		{},
	}
	return data
}

func (h *Handler) GetJobs() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			data, err := repository.GetJobs(h.DB)
			if err != nil {
				h.Logger.Errorf("Error querying job database: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			responseJson(w, data, h.Logger)
		},
	)
}

func (h *Handler) GetJobById() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			err := r.ParseForm()
			if err != nil {
				h.Logger.Errorf("Error parsing form: %v", err)
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}

			idStr := r.FormValue("id")
			id, err := strconv.Atoi(idStr)
			if err != nil {
				h.Logger.Error("Value passed as ID is invalid: %v", err)
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}
			result, err := repository.GetJobById(id, h.DB)

			responseJson(w, result, h.Logger)
		},
	)
}

func (h *Handler) GetJobByName() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			err := r.ParseForm()
			if err != nil {
				h.Logger.Errorf("Error parsing form: %v", err)
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}

			name := r.FormValue("name")
			result, err := repository.GetJobByName(name, h.DB)

			responseJson(w, result, h.Logger)
		},
	)
}

func (h *Handler) PostJob() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			err := r.ParseForm()
			if err != nil {
				h.Logger.Errorf("Error parsing form: %v", err)
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}

			job := models.Job{}

			job.Name = r.FormValue("name")
			numberStr := r.FormValue("number")
			number, err := strconv.Atoi(numberStr)
			if err != nil {
				h.Logger.Errorf("Error parsing job number value: %v", err)
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}
			job.Number = number
			job.Address = r.FormValue("address")
			job.PostCode = r.FormValue("postCode")
			job.Suburb = r.FormValue("suburb")
			job.City = r.FormValue("city")

			newJob, err := repository.PostJob(job, h.DB)
			if err != nil {
				h.Logger.Errorf("Error posting job: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			responseJson(w, newJob, h.Logger)
		},
	)
}

func (h *Handler) PutJob() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			err := r.ParseForm()
			if err != nil {
				h.Logger.Errorf("Error parsing form: %v", err)
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}

			job := models.Job{}

			idStr := r.FormValue("id")
			id, err := strconv.Atoi(idStr)
			if err != nil {
				h.Logger.Errorf("Error parsing job id value: %v", err)
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}
			job.Name = r.FormValue("name")
			numberStr := r.FormValue("number")
			number, err := strconv.Atoi(numberStr)
			if err != nil {
				h.Logger.Errorf("Error parsing job number value: %v", err)
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}
			job.Number = number
			job.Address = r.FormValue("address")
			job.PostCode = r.FormValue("postCode")
			job.Suburb = r.FormValue("suburb")
			job.City = r.FormValue("city")

			newJob, err := repository.PutJob(id, job, h.DB)
			if err != nil {
				h.Logger.Errorf("Error updating job: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			responseJson(w, newJob, h.Logger)
		},
	)
}

func (h *Handler) DeleteJob() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			err := r.ParseForm()
			if err != nil {
				h.Logger.Errorf("Error parsing form: %v", err)
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}

			idStr := r.FormValue("id")
			id, err := strconv.Atoi(idStr)
			if err != nil {
				h.Logger.Errorf("Error parsing job id value: %v", err)
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}

			deletedJob, err := repository.DeleteJob(id, h.DB)
			if err != nil {
				h.Logger.Errorf("Error deleting job: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			responseJson(w, deletedJob, h.Logger)
		},
	)
}
