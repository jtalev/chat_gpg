package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	application "github.com/jtalev/chat_gpg/application/services"
	repo "github.com/jtalev/chat_gpg/infrastructure/repository"
)

func (h *Handler) ServeIncidentReportForm() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			employeeId, err := getEmployeeId(w, r)
			if err != nil {
				log.Printf("error getting employee id: %v", err)
				http.Error(w, "error getting employee id, unauthorized", http.StatusUnauthorized)
				return
			}

			employee, err := repo.GetEmployeeByEmployeeId(employeeId, h.DB)
			if err != nil {
				log.Printf("error getting employee: %v", err)
				http.Error(w, "Not found", http.StatusNotFound)
				return
			}

			data := application.IncidentReportValues{
				ReporterId: employeeId,
				Reporter:   fmt.Sprintf("%s %s", employee.FirstName, employee.LastName),
			}

			err = executePartialTemplate(incidentReportFormPath, "incidentReportForm", data, w)
			if err != nil {
				log.Printf("error executing html template: %v", err)
				http.Error(w, "error executing html template, internal server error", http.StatusInternalServerError)
				return
			}

		},
	)
}

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

func (h *Handler) PutIncidentReportHtml() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			vals, err := parseRequestValues([]string{"uuid"}, r)
			if err != nil {
				log.Printf("error parsing request values: %v", err)
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}
			incidentReportVals, err := application.GetIncidentReport(vals[0], h.DB)
			if err != nil {
				log.Printf("error getting incident report from db: %v", err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}
			err = executePartialTemplate(putIncidentReportFormPath, "putIncidentReportForm", incidentReportVals, w)
			if err != nil {
				log.Printf("error executing html template: %v", err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}
		},
	)
}

func (h *Handler) PutIncidentReport() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			var incidentReportValues application.IncidentReportValues

			if err := json.NewDecoder(r.Body).Decode(&incidentReportValues); err != nil {
				http.Error(w, "Invalid JSON", http.StatusBadRequest)
				return
			}

			incidentReportValues, err := application.PutIncidentReport(incidentReportValues, h.DB)
			if err != nil {
				log.Printf("error updating incident report: %v", err)
				http.Error(w, "error updating incident report", http.StatusInternalServerError)
				return
			}

			err = executePartialTemplate(putIncidentReportFormPath, "putIncidentReportForm", incidentReportValues, w)
			if err != nil {
				log.Printf("error executing html template: %v", err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}
		},
	)
}

func (h *Handler) ServeSwmUserContent() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			err := executePartialTemplate(swmUserContentPath, "swmUserContent", nil, w)
			if err != nil {
				log.Printf("error executing swmUserContent.html: %v", err)
				http.Error(w, "error executing swmUserContent template", http.StatusInternalServerError)
				return
			}
		},
	)
}

func (h *Handler) GetSwmsListHtml() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			err := executePartialTemplate(swmListPath, "swmList", []int{1, 2, 3, 4}, w)
			if err != nil {
				log.Printf("error executing swmList.html: %v", err)
				http.Error(w, "error executing swmList template", http.StatusInternalServerError)
				return
			}
		},
	)
}
