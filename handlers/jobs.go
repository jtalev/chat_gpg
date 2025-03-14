package handlers

import (
	"log"
	"net/http"
	"strconv"

	application "github.com/jtalev/chat_gpg/application/services"
	"github.com/jtalev/chat_gpg/infrastructure/repository"
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
			data, err := infrastructure.GetJobs(h.DB)
			if err != nil {
				h.Logger.Errorf("Error querying job database: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			responseJSON(w, data, h.Logger)
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
			result, err := infrastructure.GetJobById(id, h.DB)

			responseJSON(w, result, h.Logger)
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
			result, err := infrastructure.GetJobByName(name, h.DB)

			responseJSON(w, result, h.Logger)
		},
	)
}

var postJobKeys = []string{"id", "name", "number", "address", "suburb", "post_code", "city"}

func (h *Handler) PostJob() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			reqVals, err := parseRequestValues(postJobKeys, r)
			if err != nil {
				log.Printf("Error parsing request values: %v", err)
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}

			jobDto := application.JobDto{
				ID:       reqVals[0],
				Name:     reqVals[1],
				Number:   reqVals[2],
				Address:  reqVals[3],
				Suburb:   reqVals[4],
				PostCode: reqVals[5],
				City:     reqVals[6],
			}

			jobDto, err = application.PostJob(jobDto, h.DB)
			if err != nil {
				h.Logger.Errorf("Error posting job: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			err = executePartialTemplate(addJobModalPath, "addJobModal", jobDto, w)
			if err != nil {
				h.Logger.Errorf("Error executing addJobModal.html: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		},
	)
}

var putJobKeys = []string{"id", "name", "number", "address", "suburb", "post_code", "city", "is_complete"}

func (h *Handler) PutJob() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			reqVals, err := parseRequestValues(putJobKeys, r)
			if err != nil {
				log.Printf("Error parsing request values: %v", err)
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}

			jobDto := application.JobDto{
				ID:         reqVals[0],
				Name:       reqVals[1],
				Number:     reqVals[2],
				Address:    reqVals[3],
				Suburb:     reqVals[4],
				PostCode:   reqVals[5],
				City:       reqVals[6],
				IsComplete: reqVals[7],
			}

			outJobDto, err := application.PutJob(jobDto, h.DB)
			if err != nil {
				log.Printf("Error updating job: %v", err)
				http.Error(w, "Not modified", http.StatusNotModified)
				return
			}

			err = executePartialTemplate(putJobModalPath, "putJobModal", outJobDto, w)
			if err != nil {
				log.Printf("Error executing template: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
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

			_, err = infrastructure.DeleteJob(id, h.DB)
			if err != nil {
				h.Logger.Errorf("Error deleting job: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			jobs, err := application.GetJobs(h.DB)
			if err != nil {
				log.Printf("Error getting job data: %v", err)
				http.Error(w, "Not found", http.StatusNotFound)
				return
			}

			err = executePartialTemplate(adminJobListPath, "adminJobList", jobs, w)
			if err != nil {
				log.Printf("Error executing adminJobList.html: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		},
	)
}
