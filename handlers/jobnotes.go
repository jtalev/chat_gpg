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
	models "github.com/jtalev/chat_gpg/domain/models"
)

func (h *Handler) initialiseJobnoteViewData() error {
	jobs, err := application.GetJobs(h.DB)
	if err != nil {
		return err
	}

	activeJobs := make([]models.Job, 0)
	for _, job := range jobs {
		if !job.IsComplete {
			activeJobs = append(activeJobs, job)
		}
	}

	err = h.jobnotes.InitialJobnoteViewData(activeJobs)
	if err != nil {
		return err
	}
	return nil
}

func (h *Handler) ServeJobsView() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			err := h.initialiseJobnoteViewData()
			if err != nil {
				log.Printf("error initialising view data: %v", err)
				http.Error(w, "error initialising view data", http.StatusInternalServerError)
				return
			}

			component := "jobs"
			title := "Jobs - GPG"
			renderTemplate(w, r, component, title, h.jobnotes.JobnoteViewData)
		},
	)
}

func servePaintNoteForm(path string, data jobnotes.PaintnoteFormData, w http.ResponseWriter) error {
	err := executePartialTemplate(path, "paintNoteForm", data, w)
	if err != nil {
		log.Printf("error executing template paintNoteForm: %v", err)
		http.Error(w, "error executing template, internal server error", http.StatusInternalServerError)
		return err
	}
	return nil
}

func serveTaskNoteForm(path string, data jobnotes.TasknoteFormData, w http.ResponseWriter) error {
	err := executePartialTemplate(path, "taskNoteForm", data, w)
	if err != nil {
		log.Printf("error executing template taskNoteForm: %v", err)
		http.Error(w, "error executing template, internal server error", http.StatusInternalServerError)
		return err
	}
	return nil
}

func serveImageNoteForm(path string, data jobnotes.ImagenoteFormData, w http.ResponseWriter) error {
	err := executePartialTemplate(path, "imageNoteForm", data, w)
	if err != nil {
		log.Printf("error executing template imageNoteForm: %v", err)
		http.Error(w, "error executing template, internal server error", http.StatusInternalServerError)
		return err
	}
	return nil
}

func (h *Handler) ServeNoteForm() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			reqVals, err := parseRequestValues([]string{"uuid", "note_type", "job_id"}, r)
			if err != nil {
				log.Printf("error parsing request values: %v", err)
				http.Error(w, "error parsing request values, bad request", http.StatusBadRequest)
				return
			}

			uuid, noteType := reqVals[0], reqVals[1]
			jobIdStr := reqVals[2]
			jobId, err := strconv.Atoi(jobIdStr)
			if err != nil {
				log.Printf("error converting job ID to int: %v", err)
				http.Error(w, "error converting job ID to int, bad request", http.StatusBadRequest)
				return
			}

			if uuid == "" {
				path := ""
				switch noteType {
				case "paint_note":
					path = paintNoteFormPath
					h.jobnotes.PaintnoteFormData.FormType = "post"
					h.jobnotes.PaintnoteFormData.JobId = jobId
					if err := servePaintNoteForm(path, h.jobnotes.PaintnoteFormData, w); err != nil {
						return
					}
				case "task_note":
					path = taskNoteFormPath
					h.jobnotes.TasknoteFormData.FormType = "post"
					h.jobnotes.TasknoteFormData.JobId = jobId
					if err := serveTaskNoteForm(path, h.jobnotes.TasknoteFormData, w); err != nil {
						return
					}
				case "image_note":
					path = imageNoteFormPath
					h.jobnotes.ImagenoteFormData.FormType = "post"
					h.jobnotes.ImagenoteFormData.JobId = jobId
					if err := serveImageNoteForm(path, h.jobnotes.ImagenoteFormData, w); err != nil {
						return
					}
				default:
					log.Printf("note type %s not supported", noteType)
					http.Error(w, "note type not supported, bad request", http.StatusBadRequest)
					return
				}
			} else {

			}
		},
	)
}

func (h *Handler) ServeJobnoteTiles() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			err := h.initialiseJobnoteViewData()
			if err != nil {
				log.Printf("error initialising view data: %v", err)
				http.Error(w, "error initialising view data", http.StatusInternalServerError)
				return
			}

			err = executePartialTemplate(jobnoteTilesPath, "jobnoteTiles", h.jobnotes.JobnoteViewData, w)
			if err != nil {
				log.Printf("error executing template: %v", err)
				http.Error(w, "error executing template, internal server error", http.StatusInternalServerError)
				return
			}
		},
	)
}

func (h *Handler) GetJobNotes() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			reqVals, err := parseRequestValues([]string{"job_id"}, r)
			if err != nil {
				log.Printf("error parsing request values: %v", err)
				http.Error(w, "error parsing request values, bad request", http.StatusBadRequest)
				return
			}

			jobid, err := strconv.Atoi(reqVals[0])
			if err != nil {
				log.Printf("error converting request value to int: %v", err)
				http.Error(w, "error converting request value to int, bad request", http.StatusBadRequest)
				return
			}

			err = h.jobnotes.GetJobNotes(jobid)
			if err != nil {
				log.Printf("error getting notes by job: %v", err)
				http.Error(w, "error getting notes by job, bad request", http.StatusBadRequest)
				return
			}

			job, err := application.GetJobById(jobid, h.DB)
			if err != nil {
				log.Printf("error getting job: %v", err)
				http.Error(w, "error getting job", http.StatusInternalServerError)
				return
			}

			h.jobnotes.JobSummary.ID = job.ID
			h.jobnotes.JobSummary.Name = job.Name
			h.jobnotes.JobSummary.Address = fmt.Sprintf("%d %s", job.Number, job.Address)
			h.jobnotes.JobSummary.PaintnoteCount = len(h.jobnotes.Paintnotes)
			h.jobnotes.JobSummary.TasknoteCount = len(h.jobnotes.Tasknotes)
			h.jobnotes.JobSummary.ImagenoteCount = len(h.jobnotes.Imagenotes)

			templatePaths := []string{jobNotesPath, paintNotePath, taskNotePath, imageNotePath}
			err = h.ServeMultiTemplate(templatePaths, "jobNotes", h.jobnotes, w)
			if err != nil {
				log.Printf("error serving html: %v", err)
				return
			}
		},
	)
}

func unmarshalNote(requestBody []byte, j *jobnotes.Jobnotes) error {
	log.Printf("jobnotes: %v", j.Note)
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

			switch h.jobnotes.Note.NoteType {
			case "paint_note":
				if err := servePaintNoteForm(paintNoteFormPath, h.jobnotes.PaintnoteFormData, w); err != nil {
					return
				}
			case "task_note":
				if err := serveTaskNoteForm(taskNoteFormPath, h.jobnotes.TasknoteFormData, w); err != nil {
					return
				}
			case "image_note":
				if err := serveImageNoteForm(imageNoteFormPath, h.jobnotes.ImagenoteFormData, w); err != nil {
					return
				}
			default:
				log.Printf("note type %s not supported", h.jobnotes.Note.NoteType)
				http.Error(w, "note type not supported, bad request", http.StatusBadRequest)
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
