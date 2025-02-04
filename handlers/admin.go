package handlers

import (
	"html/template"
	"log"
	"net/http"

	infrastructure "github.com/jtalev/chat_gpg/infrastructure/repository"
)

type AdminData struct {
}

func getAdminData() []AdminData {
	data := []AdminData{
		{},
	}
	return data
}

func (h *Handler) RenderJobTab() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			data, err := infrastructure.GetJobs(h.DB)
			if err != nil {
				log.Printf("Error querying job database: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			tmpl, err := template.ParseFiles(
				adminJobTabPath,
				adminJobListPath,
			)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			err = tmpl.ExecuteTemplate(w, "adminJobTab", data)
			if err != nil {
				log.Printf("Error executing template: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		},
	)
}

func (h *Handler) RenderJobList() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			data, err := infrastructure.GetJobs(h.DB)
			if err != nil {
				log.Printf("Error querying job database: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			err = executePartialTemplate(adminJobListPath, "adminJobList", data, w)
			if err != nil {
				log.Printf("Error executing adminJobList.html: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		},
	)
}
