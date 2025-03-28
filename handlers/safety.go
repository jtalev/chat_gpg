package handlers

import (
	"log"
	"net/http"

	application "github.com/jtalev/chat_gpg/application/services"
)

type IncidentReportValues struct {
	FullName            string `json:"full_name"`
	HomeAddress         string `json:"home_address"`
	ContactNumber       string `json:"contact_number"`
	IncidentDate        string `json:"incident_date"`
	Time                string `json:"time"`
	PoliceNotified      string `json:"police_notified"`
	IncidentLocation    string `json:"incident_location"`
	IncidentDescription string `json:"incident_description"`
	WasWitnessed        string `json:"was_witnessed"`
	WasInjured          string `json:"was_injured"`
	FurtherDetails      string `json:"further_details"`
	WasTreated          string `json:"was_treated"`
	TreatmentLocation   string `json:"treatment_location"`
	IncInfoDate1        string `json:"inc_info_date_1"`
	IncInfoDate2        string `json:"inc_info_date_2"`
	IncInfoDate3        string `json:"inc_info_date_3"`
	IncInfoDate4        string `json:"inc_info_date_4"`
	IncInfoDate5        string `json:"inc_info_date_5"`
	IncInfoAction1      string `json:"inc_info_action_1"`
	IncInfoAction2      string `json:"inc_info_action_2"`
	IncInfoAction3      string `json:"inc_info_action_3"`
	IncInfoAction4      string `json:"inc_info_action_4"`
	IncInfoAction5      string `json:"inc_info_action_5"`
	IncInfoName1        string `json:"inc_info_name_1"`
	IncInfoName2        string `json:"inc_info_name_2"`
	IncInfoName3        string `json:"inc_info_name_3"`
	IncInfoName4        string `json:"inc_info_name_4"`
	IncInfoName5        string `json:"inc_info_name_5"`
	Reporter            string `json:"reporter"`
	Signature           string `json:"signature"`
	ReportDate          string `json:"report_date"`
	// FullNameErr            string
	// HomeAddressErr         string
	// ContactNumberErr       string
	// DateErr                string
	// PoliceNotifiedErr      string
	// IncidentLocationErr    string
	// IncidentDescriptionErr string
	// WasWitnessedErr        string
	// VictimInjuredErr       string
	// TreatmentProvidedErr   string
	// TreatmentLocationErr   string
	// ReporterErr            string
	// SignatureErr           string
	// ReportDateErr          string
}

var incidentReportValues = IncidentReportValues{
	FullName:            "Josh Talev",
	HomeAddress:         "95 Jedda St",
	ContactNumber:       "0450 579 387",
	IncidentDate:        "27-03-2025",
	Time:                "10:30am",
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
	Reporter:            "Hayden Glenny",
	ReportDate:          "27-03-2025",
}

var incidentReportPath = "../ui/static/files/incident_report_forms.pdf"
var incidentReportDestPath = "../ui/static/files/incident_report_test.pdf"
var tempJsonPath = "../ui/static/files/incident_form_data.json"
var exportedJsonPath = "../ui/static/files/exported_form_data.json"

func (h *Handler) GenerateIncidentReport() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			err := application.GeneratePdf(incidentReportPath, tempJsonPath, incidentReportDestPath, incidentReportValues)
			if err != nil {
				log.Printf("Error executing json template: %v", err)
			}
		},
	)
}
