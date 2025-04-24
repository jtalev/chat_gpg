package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	application "github.com/jtalev/chat_gpg/application/services"
	"github.com/jtalev/chat_gpg/application/services/safety"
	models "github.com/jtalev/chat_gpg/domain/models"
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

			data := safety.IncidentReportValues{
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
			var incidentReportValues safety.IncidentReportValues

			if err := json.NewDecoder(r.Body).Decode(&incidentReportValues); err != nil {
				http.Error(w, "Invalid JSON", http.StatusBadRequest)
				return
			}

			incidentReportValues, err := safety.GenerateIncidentReportPdf(incidentReportValues, h.DB)
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

			pdfUrl, err := safety.GetIncidentReportUrl(vals[0], h.DB)
			if err != nil {
				log.Printf("error getting incident report url: %v", err)
				http.Error(w, "error getting incident report url", http.StatusInternalServerError)
				return
			}

			data := struct {
				URL   string
				Scale float32
			}{
				URL:   pdfUrl,
				Scale: 1.3,
			}

			err = executePartialTemplate(iframePdfPath, "iframePdf", data, w)
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

			err = safety.DeleteIncidentReport(vals[0], h.DB)
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
			incidentReportVals, err := safety.GetIncidentReport(vals[0], h.DB)
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
			var incidentReportValues safety.IncidentReportValues

			if err := json.NewDecoder(r.Body).Decode(&incidentReportValues); err != nil {
				http.Error(w, "Invalid JSON", http.StatusBadRequest)
				return
			}

			incidentReportValues, err := safety.PutIncidentReport(incidentReportValues, h.DB)
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

var s = safety.Swms{}

var defaultSwm = safety.Swms{
	Swms: models.Swms{

		HighRiskWorks:          true,
		SafetyBoots:            true,
		HiVisClothing:          true,
		StepBelow2m:            true,
		Step1:                  "Site Set Up",
		Hazards1:               "Muscle strains/ additional tripping hazards",
		Risks1:                 "Muscle strain and sprains from incorrect lifting technique/ tripping hazard over site storage position",
		InitialRisk1:           "B3",
		ControlMeasures1:       "Warm Up stretches, correct lifting technique, 2 person lifts when required/ preplanning to find suitable storage area for equipment",
		ResidualRisk1:          "A2",
		ControlResponsibility1: "Project Manager Foreman Painters",
		Step2:                  "Working on Ladders",
		Hazards2:               "Fall Hazard",
		Risks2:                 "Potential serious injury from falling from height of ladder",
		InitialRisk2:           "D3",
		ControlMeasures2:       "Use of platform steps whenever possible, no overreaching of ladder, where possible utilise mobile equipment or plant",
		ResidualRisk2:          "B2",
		ControlResponsibility2: "Project Manager Foreman Painters",
		Step3:                  "Working in Public Spaces",
		Hazards3:               "Public Safety",
		Risks3:                 "Risk of injury or harm to the general public",
		InitialRisk3:           "C2",
		ControlMeasures3:       "Isolate working area with the use of bollards, hard barriers and signage",
		ResidualRisk3:          "B1",
		ControlResponsibility3: "Project Manager Foreman Painters",
		Step4:                  "General Preparation works",
		Hazards4:               "Flying debris from sanding/ respiratory damage from sanding of paint and plaster",
		Risks4:                 "Eye damage from flying debris while undertaking preparation works/ Respiratory damage from paint/plaster dust",
		InitialRisk4:           "B4",
		ControlMeasures4:       "Use of correct PPE while undertaking preparation works including safety glasses/ P2 rated mask if required",
		ResidualRisk4:          "B1",
		ControlResponsibility4: "Project Manager Foreman Painters",
		Step5:                  "General Painting Works",
		Hazards5:               "Muscle strains, eye damage, fume inhalation",
		Risks5:                 "Muscle strains from overuse/ flying debris when working overhead/ fume inhalation",
		InitialRisk5:           "B4",
		ControlMeasures5:       "Take breaks when required and continue stretches throughout the day, use of safety glasses working overhead/ P2 mask when required",
		ResidualRisk5:          "B2",
		ControlResponsibility5: "Project Manager Foreman Painters",
	},
}

func (h *Handler) GetSwmsListHtml() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			s.Db = h.DB
			employeeIdVal := r.Context().Value("employee_id")
			employeeId, ok := employeeIdVal.(string)
			if !ok {
				log.Printf("error asserting employee_id from context: %v", employeeIdVal)
				http.Error(w, "unauthorized or invalid employee id", http.StatusBadRequest)
				return
			}

			err := s.GetUserRole(employeeId)
			if err != nil {
				log.Printf("error getting user role: %v", err)
				http.Error(w, "error getting user role, unauthorized", http.StatusUnauthorized)
				return
			}

			err = s.GetSwms()
			if err != nil {
				log.Printf("error getting swms: %v", err)
				http.Error(w, "error getting swms, internal server error", http.StatusInternalServerError)
				return
			}

			err = executePartialTemplate(swmListPath, "swmList", s, w)
			if err != nil {
				log.Printf("error executing swmList.html: %v", err)
				http.Error(w, "error executing swmList template", http.StatusInternalServerError)
				return
			}
		},
	)
}

func (h *Handler) GenerateSwmsPdf() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			err := json.NewDecoder(r.Body).Decode(&s.Swms)
			if err != nil {
				log.Printf("error decoding swms json data: %v", err)
				http.Error(w, "error decoding swms json data, bad request", http.StatusBadRequest)
				return
			}

			s.Db = h.DB
			s.Errors, err = s.PostSwm(s.Swms)
			if err != nil {
				log.Printf("error posting swms: %v", err)
				http.Error(w, "error posting swms, internal server error", http.StatusInternalServerError)
				return
			}
			s.GenerateSwmsPdf(s.Swms)
		},
	)
}

func (h *Handler) ServeSwmsPdf() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			uuid, err := parseRequestValues([]string{"uuid"}, r)
			if err != nil {
				log.Printf("error parsing request values: %v", err)
				http.Error(w, "error parsing request values, bad request", http.StatusBadRequest)
				return
			}

			s.Db = h.DB
			pdfUrl, err := s.GetSwmsPdfUrl(uuid[0])
			if err != nil {
				log.Printf("error getting swms url: %v", err)
				http.Error(w, "error getting swms url, internal server error", http.StatusInternalServerError)
				return
			}

			data := struct {
				URL   string
				Scale float32
			}{
				URL:   pdfUrl,
				Scale: 1,
			}

			err = executePartialTemplate(iframePdfPath, "iframePdf", data, w)
			if err != nil {
				log.Printf("error executing iframePdf template: %v", err)
				http.Error(w, "error executing iframePdf template, internal server error", http.StatusInternalServerError)
				return
			}
		},
	)
}

func (h *Handler) ServeSwmsForm() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			s.Errors = models.SwmsErrors{}
			s.Swms = defaultSwm.Swms
			jobs, err := application.GetJobs(h.DB)
			if err != nil {
				log.Printf("error fetching jobs: %v", err)
				http.Error(w, "error fetching jobs, internal server error", http.StatusInternalServerError)
				return
			}
			s.Jobs = jobs
			err = executePartialTemplate(swmsFormPath, "swmsForm", s, w)
			if err != nil {
				log.Printf("error executing swmsForm template: %v", err)
				http.Error(w, "error executing swmsForm template, internal server error", http.StatusInternalServerError)
				return
			}
		},
	)
}

func (h *Handler) ServeSwmsFormPut() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			log.Println("serve swms form put")
			s.Errors = models.SwmsErrors{}
			s.Db = h.DB
			uuid, err := parseRequestValues([]string{"uuid"}, r)
			if err != nil {
				log.Printf("error parsing request values: %v", err)
				http.Error(w, "error parsing request values, bad request", http.StatusBadRequest)
				return
			}
			s.Swms.UUID = uuid[0]
			err = s.GetSwmsByUUID()
			if err != nil {
				log.Printf("error fetching swms by UUID: %v", err)
				http.Error(w, "error fetching swms by UUID, internal server error", http.StatusInternalServerError)
				return
			}

			jobs, err := application.GetJobs(h.DB)
			if err != nil {
				log.Printf("error fetching jobs: %v", err)
				http.Error(w, "error fetching jobs, internal server error", http.StatusInternalServerError)
				return
			}
			s.Jobs = jobs
			err = executePartialTemplate(updateSwmsFormPath, "updateSwmsForm", s, w)
			if err != nil {
				log.Printf("error executing updateSwmsForm template: %v", err)
				http.Error(w, "error executing updateSwmsForm template, internal server error", http.StatusInternalServerError)
				return
			}
		},
	)
}

var reqParams = []string{"uuid", "job_id", "project_activity", "project_number", "site_address", "contact_name", "contact_number",
	"email_address", "swms_date", "high_risk_works", "safety_gloves", "safety_boots", "safety_glasses",
	"protective_clothing", "respiratory_protection", "hi_vis_clothing", "safety_helmet", "fall_arrest",
	"other_1", "other_2", "step_below_2m", "step_above_2m", "scaffold", "pressure_washer_diesel",
	"roof_anchor_points", "extension_ladder", "electric_scissor_lift", "diesel_scissor_lift",
	"electric_knuckle_boom", "diesel_knuckle_boom", "airless_spray_gun", "angle_grinder",
	"step_1", "hazards_1", "risks_1", "initial_risk_1", "control_measures_1", "residual_risk_1", "control_responsibility_1",
	"step_2", "hazards_2", "risks_2", "initial_risk_2", "control_measures_2", "residual_risk_2", "control_responsibility_2",
	"step_3", "hazards_3", "risks_3", "initial_risk_3", "control_measures_3", "residual_risk_3", "control_responsibility_3",
	"step_4", "hazards_4", "risks_4", "initial_risk_4", "control_measures_4", "residual_risk_4", "control_responsibility_4",
	"step_5", "hazards_5", "risks_5", "initial_risk_5", "control_measures_5", "residual_risk_5", "control_responsibility_5",
	"step_6", "hazards_6", "risks_6", "initial_risk_6", "control_measures_6", "residual_risk_6", "control_responsibility_6",
	"step_7", "hazards_7", "risks_7", "initial_risk_7", "control_measures_7", "residual_risk_7", "control_responsibility_7",
	"step_8", "hazards_8", "risks_8", "initial_risk_8", "control_measures_8", "residual_risk_8", "control_responsibility_8",
	"step_9", "hazards_9", "risks_9", "initial_risk_9", "control_measures_9", "residual_risk_9", "control_responsibility_9",
	"step_10", "hazards_10", "risks_10", "initial_risk_10", "control_measures_10", "residual_risk_10", "control_responsibility_10",
	"step_11", "hazards_11", "risks_11", "initial_risk_11", "control_measures_11", "residual_risk_11", "control_responsibility_11",
	"step_12", "hazards_12", "risks_12", "initial_risk_12", "control_measures_12", "residual_risk_12", "control_responsibility_12",
	"sign_date_1", "sign_name_1", "sign_sig_1",
	"sign_date_2", "sign_name_2", "sign_sig_2",
	"sign_date_3", "sign_name_3", "sign_sig_3",
	"sign_date_4", "sign_name_4", "sign_sig_4",
	"sign_date_5", "sign_name_5", "sign_sig_5",
	"sign_date_6", "sign_name_6", "sign_sig_6",
	"sign_date_7", "sign_name_7", "sign_sig_7",
	"sign_date_8", "sign_name_8", "sign_sig_8",
	"sign_date_9", "sign_name_9", "sign_sig_9"}

var swmsMap = map[string]*bool{
	"high_risk_works":        &s.Swms.HighRiskWorks,
	"safety_gloves":          &s.Swms.SafetyGloves,
	"safety_boots":           &s.Swms.SafetyBoots,
	"safety_glasses":         &s.Swms.SafetyGlasses,
	"protective_clothing":    &s.Swms.ProtectiveClothing,
	"respiratory_protection": &s.Swms.RespiratoryProtection,
	"hi_vis_clothing":        &s.Swms.HiVisClothing,
	"safety_helmet":          &s.Swms.SafetyHelmet,
	"fall_arrest":            &s.Swms.FallArrest,
	"step_below_2m":          &s.Swms.StepBelow2m,
	"step_above_2m":          &s.Swms.StepAbove2m,
	"scaffold":               &s.Swms.Scaffold,
	"pressure_washer_diesel": &s.Swms.PressureWasherDiesel,
	"roof_anchor_points":     &s.Swms.RoofAnchorPoints,
	"extension_ladder":       &s.Swms.ExtensionLadder,
	"electric_scissor_lift":  &s.Swms.ElectricScissorLift,
	"diesel_scissor_lift":    &s.Swms.DieselScissorLift,
	"electric_knuckle_boom":  &s.Swms.ElectricKnuckleBoom,
	"diesel_knuckle_boom":    &s.Swms.DieselKnuckleBoom,
	"airless_spray_gun":      &s.Swms.AirlessSprayGun,
	"angle_grinder":          &s.Swms.AngleGrinder,
}

func (h *Handler) PostSwms() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			reqVals, err := parseRequestValues(reqParams, r)
			if err != nil {
				log.Printf("error parsing request values: %v", err)
				http.Error(w, "error parsing request values, bad request", http.StatusBadRequest)
				return
			}

			for i, param := range reqParams {
				val := reqVals[i]
				if val == "true" {
					if swmsField, ok := swmsMap[param]; ok {
						*swmsField = true
						continue
					}
				}

				swmsVal := reflect.ValueOf(&s.Swms).Elem()
				swmsType := swmsVal.Type()

				for j := 0; j < swmsType.NumField(); j++ {
					field := swmsType.Field(j)
					fieldName := field.Tag.Get("json")
					if fieldName == "" {
						fieldName = strings.ToLower(field.Name)
					}

					if fieldName == param {
						fieldVal := swmsVal.FieldByName(field.Name)
						if fieldVal.CanSet() {
							switch fieldVal.Kind() {
							case reflect.String:
								fieldVal.SetString(val)
							case reflect.Int:
								log.Println(val)
								valInt, err := strconv.Atoi(val)
								if err != nil {
									log.Printf("error converting string to int: %v", err)
									http.Error(w, "error converting string to int", http.StatusConflict)
									return
								}
								fieldVal.SetInt(int64(valInt))
							default:
								log.Printf("Unsupported field type: %s", fieldVal.Kind())
							}
						}
					}
				}
			}
			jobs, err := application.GetJobs(h.DB)
			if err != nil {
				log.Printf("error fetching jobs: %v", err)
				http.Error(w, "error fetching jobs, internal server error", http.StatusInternalServerError)
				return
			}
			s.Jobs = jobs

			s.Errors, err = s.PostSwm(s.Swms)
			if err != nil {
				log.Printf("error posting swms: %v", err)
				http.Error(w, "error posting swms, internal server error", http.StatusInternalServerError)
				return
			}
			if s.Errors.IsSuccessful {
				s.GenerateSwmsPdf(s.Swms)
			}

			err = executePartialTemplate(swmsFormPath, "swmsForm", s, w)
			if err != nil {
				log.Printf("error executing swmsForm: %v", err)
				http.Error(w, "error executing swmsForm, internal server error", http.StatusInternalServerError)
				return
			}
		},
	)
}

func (h *Handler) PutSwms() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			s.Db = h.DB
			reqVals, err := parseRequestValues(reqParams, r)
			if err != nil {
				log.Printf("error parsing request values: %v", err)
				http.Error(w, "error parsing request values, bad request", http.StatusBadRequest)
				return
			}

			for i, param := range reqParams {
				val := reqVals[i]
				if val == "true" {
					if swmsField, ok := swmsMap[param]; ok {
						*swmsField = true
						continue
					}
				}

				swmsVal := reflect.ValueOf(&s.Swms).Elem()
				swmsType := swmsVal.Type()

				for j := 0; j < swmsType.NumField(); j++ {
					field := swmsType.Field(j)
					fieldName := field.Tag.Get("json")
					if fieldName == "" {
						fieldName = strings.ToLower(field.Name)
					}

					if fieldName == param {
						fieldVal := swmsVal.FieldByName(field.Name)
						if fieldVal.CanSet() {
							switch fieldVal.Kind() {
							case reflect.String:
								fieldVal.SetString(val)
							case reflect.Int:
								log.Println(val)
								valInt, err := strconv.Atoi(val)
								if err != nil {
									log.Printf("error converting string to int: %v", err)
									http.Error(w, "error converting string to int", http.StatusConflict)
									return
								}
								fieldVal.SetInt(int64(valInt))
							default:
								log.Printf("Unsupported field type: %s", fieldVal.Kind())
							}
						}
					}
				}
			}
			jobs, err := application.GetJobs(h.DB)
			if err != nil {
				log.Printf("error fetching jobs: %v", err)
				http.Error(w, "error fetching jobs, internal server error", http.StatusInternalServerError)
				return
			}
			s.Jobs = jobs

			s.Errors, err = s.PutSwms(s.Swms)
			if err != nil {
				log.Printf("error posting swms: %v", err)
				http.Error(w, "error posting swms, internal server error", http.StatusInternalServerError)
				return
			}
			if s.Errors.IsSuccessful {
				s.GenerateSwmsPdf(s.Swms)
			}

			err = executePartialTemplate(swmsFormPath, "swmsForm", s, w)
			if err != nil {
				log.Printf("error executing swmsForm: %v", err)
				http.Error(w, "error executing swmsForm, internal server error", http.StatusInternalServerError)
				return
			}
		},
	)
}

func (h *Handler) DeleteSwms() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			vals, err := parseRequestValues([]string{"uuid"}, r)
			if err != nil {
				log.Printf("error parsing request values: %v", err)
				http.Error(w, "error parsing request values, bad request", http.StatusBadRequest)
				return
			}

			s.Db = h.DB
			s.Swms.UUID = vals[0]
			s.DeleteSwms()
		},
	)
}
