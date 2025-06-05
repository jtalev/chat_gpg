package handlers

import (
	"log"
	"net/http"
	"strconv"
)

func (h *Handler) GetJobNotes() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			_jobid, err := parseRequestValues([]string{"uuid"}, r)
			if err != nil {
				log.Printf("error parsing request values: %v", err)
				http.Error(w, "error parsing request values, bad request", http.StatusBadRequest)
				return
			}

			jobid, err := strconv.Atoi(_jobid[0])
			if err != nil {
				log.Printf("error converting request value to int: %v", err)
				http.Error(w, "error converting request value to in, bad request", http.StatusBadRequest)
				return
			}

			err = h.jobnotes.GetJobNotes(jobid)
			if err != nil {
				log.Printf("error getting notes by job: %v", err)
				http.Error(w, "error getting notes by job, bad request", http.StatusBadRequest)
				return
			}

			log.Printf("paint notes: %v", h.jobnotes.Paintnotes)
			log.Printf("task notes: %v", h.jobnotes.Tasknotes)
			log.Printf("image notes: %v", h.jobnotes.Imagenotes)
		},
	)
}
