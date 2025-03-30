package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	application "github.com/jtalev/chat_gpg/application/services"
)

func (h *Handler) GenerateIncidentReport() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			var incidentReportValues application.IncidentReportValues

			if err := json.NewDecoder(r.Body).Decode(&incidentReportValues); err != nil {
				http.Error(w, "Invalid JSON", http.StatusBadRequest)
				return
			}

			incidentReportValues, err := application.GenerateIncidentReportPdf(incidentReportValues, h.DB)
			if err != nil {
				log.Printf("error generating incident report pdf: %v", err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}

			err = executePartialTemplate(incidentReportFormPath, "incidentReportForm", incidentReportValues, w)
			if err != nil {
				log.Printf("error executing incident report form template: %v", err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}
		},
	)
}

func (h *Handler) GetIncidentReport() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			vals, err := parseRequestValues([]string{"uuid"}, r)
			if err != nil {
				log.Printf("error parsing uuid from request: %v", err)
				http.Error(w, "error parsing uuid from request", http.StatusBadRequest)
				return
			}

			pdfUrl, err := application.GetIncidentReportUrl(vals[0], h.DB)
			if err != nil {
				log.Printf("error getting incident report url: %v", err)
				http.Error(w, "error getting incident report url", http.StatusInternalServerError)
				return
			}

			err = executePartialTemplate(iframePdfPath, "iframePdf", pdfUrl, w)
			if err != nil {
				log.Printf("error executing incident report form template: %v", err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}
		},
	)
}

func (h *Handler) DeleteIncidentReport() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			vals, err := parseRequestValues([]string{"uuid"}, r)
			if err != nil {
				log.Printf("error parsing uuid from request: %v", err)
				http.Error(w, "error parsing uuid from request", http.StatusBadRequest)
				return
			}

			err = application.DeleteIncidentReport(vals[0], h.DB)
			if err != nil {
				log.Printf("error deleting incident report url: %v", err)
				http.Error(w, "error deleting incident report url", http.StatusInternalServerError)
				return
			}

			fmt.Fprint(w, "View PDF here.")
		},
	)
}
