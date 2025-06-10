package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	application "github.com/jtalev/chat_gpg/application/services"
	"github.com/jtalev/chat_gpg/application/services/jobnotes"
)

type jobnoteSummary struct {
	Name           string
	Address        string
	PaintnoteCount int
	TasknoteCount  int
	ImagenoteCount int
}

type jobnoteViewData struct {
	JobCount  int
	Summaries []jobnoteSummary
}

func (h *Handler) initialiseJobnoteViewData() (jobnoteViewData, error) {
	var data jobnoteViewData
	jobs, err := application.GetJobs(h.DB)
	if err != nil {
		return data, err
	}
	data.JobCount = len(jobs)

	for _, j := range jobs {
		data.Summaries = append(data.Summaries, jobnoteSummary{
			Name:    j.Name,
			Address: fmt.Sprintf("%d %s", j.Number, j.Address),
		})
	}
	return data, nil
}

func (h *Handler) ServeJobsView() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {

			data, err := h.initialiseJobnoteViewData()
			if err != nil {
				log.Printf("error initialising view data: %v", err)
				http.Error(w, "error initialising view data", http.StatusInternalServerError)
				return
			}

			component := "jobs"
			title := "Jobs - GPG"
			renderTemplate(w, r, component, title, data)
		},
	)
}

func (h *Handler) GetJobNotes() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			reqVals, err := parseRequestValues([]string{"uuid"}, r)
			if err != nil {
				log.Printf("error parsing request values: %v", err)
				http.Error(w, "error parsing request values, bad request", http.StatusBadRequest)
				return
			}

			jobid, err := strconv.Atoi(reqVals[0])
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

func unmarshalNote(requestBody []byte, j *jobnotes.Jobnotes) error {
	if err := json.Unmarshal(requestBody, &j.Note); err != nil {
		return err
	}

	return nil
}

func unmarshalNotetype(requestBody []byte, j *jobnotes.Jobnotes) error {
	switch j.Note.NoteType {
	case "paint_note":
		if err := json.Unmarshal(requestBody, &j.Paintnote); err != nil {
			return err
		}
	case "task_note":
		if err := json.Unmarshal(requestBody, &j.Tasknote); err != nil {
			return err
		}
	case "image_note":
		if err := json.Unmarshal(requestBody, &j.Imagenote); err != nil {
			return err
		}
	default:
		err := fmt.Errorf("note type %s not supported", j.Note.NoteType)
		return err
	}
	return nil
}

func (h *Handler) PostNote() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			requestBody, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "failed to read body", http.StatusInternalServerError)
				return
			}

			err = unmarshalNote(requestBody, h.jobnotes)
			if err != nil {
				log.Printf("bad json, error unmarshaling note: %v", err)
				http.Error(w, "bad json, error unmarshaling note", http.StatusBadRequest)
				return
			}

			err = unmarshalNotetype(requestBody, h.jobnotes)
			if err != nil {
				log.Printf("bad json, error unmarshaling note type: %v", err)
				http.Error(w, "bad json, error unmarshaling note type", http.StatusBadRequest)
				return
			}

			err = h.jobnotes.PostNote(h.jobnotes.Note.NoteType)
			if err != nil {
				log.Printf("error posting note: %v", err)
				http.Error(w, "error posting note", http.StatusInternalServerError)
				return
			}
		},
	)
}

func (h *Handler) PutNote() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			requestBody, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "failed to read body", http.StatusInternalServerError)
				return
			}

			err = unmarshalNote(requestBody, h.jobnotes)
			if err != nil {
				log.Printf("bad json, error unmarshaling note: %v", err)
				http.Error(w, "bad json, error unmarshaling note", http.StatusBadRequest)
				return
			}

			err = unmarshalNotetype(requestBody, h.jobnotes)
			if err != nil {
				log.Printf("bad json, error unmarshaling note type: %v", err)
				http.Error(w, "bad json, error unmarshaling note type", http.StatusBadRequest)
				return
			}

			err = h.jobnotes.PutNote(h.jobnotes.Note.NoteType)
			if err != nil {
				log.Printf("error udating note: %v", err)
				http.Error(w, "error updating note", http.StatusInternalServerError)
				return
			}
		},
	)
}

func (h *Handler) DeleteNote() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if ok := h.DecodeJson(&h.jobnotes.Note, w, r); !ok {
				return
			}

			err := h.jobnotes.DeleteNote(h.jobnotes.Note.Uuid)
			if err != nil {
				log.Printf("error deleting note: %v", err)
				http.Error(w, "error deleting note", http.StatusInternalServerError)
				return
			}
		},
	)
}
