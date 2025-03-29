package handlers

import (
	"log"
	"net/http"

	application "github.com/jtalev/chat_gpg/application/services"
)

var incidentReportValues = application.IncidentReportValues{
	ReporterId:          "3233123",
	FullName:            "Josh Talev",
	HomeAddress:         "95 Jedda St",
	ContactNumber:       "0450 579 387",
	IncidentDate:        "27-03-2025",
	IncidentTime:        "10:30am",
	PoliceNotified:      "yes",
	IncidentLocation:    "work",
	IncidentDescription: "knocked out",
	WasWitnessed:        "no",
	WasInjured:          "yes",
	FurtherDetails:      "",
	WasTreated:          "refused",
	TreatmentLocation:   "onsite",
	IncInfoDate1:        "",
	IncInfoDate2:        "",
	IncInfoDate3:        "",
	IncInfoDate4:        "",
	IncInfoDate5:        "",
	IncInfoAction1:      "",
	IncInfoAction2:      "",
	IncInfoAction3:      "",
	IncInfoAction4:      "",
	IncInfoAction5:      "",
	IncInfoName1:        "",
	IncInfoName2:        "",
	IncInfoName3:        "",
	IncInfoName4:        "",
	IncInfoName5:        "",
	Reporter:            "Russell Graham",
	Signature:           "RusselGraham",
	ReportDate:          "27-03-2025",
}

func (h *Handler) GenerateIncidentReport() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			err := application.GenerateIncidentReportPdf(incidentReportValues, h.DB)
			if err != nil {
				log.Printf("error generating incident report pdf: %v", err)
				http.Error(w, "internal server errror", http.StatusInternalServerError)
				return
			}
		},
	)
}
