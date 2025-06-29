package handlers

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jtalev/chat_gpg/application/services"
	"github.com/jtalev/chat_gpg/domain/models"
	"github.com/jtalev/chat_gpg/infrastructure/repository"
	"go.uber.org/zap"
)

func (h *Handler) RenderLeaveHistoryTab() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			employeeId, err := getEmployeeId(w, r)
			if err != nil {
				log.Println("Error retrieving employee_id:", err)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			historyData, err := application.GetLeaveHistoryByEmployeeId(employeeId, h.DB)
			if err != nil {
				log.Println("Error retrieving leave history from db:", err)
				http.Error(w, "Status not found", http.StatusNotFound)
				return
			}

			tmpl, err := template.ParseFiles(leaveHistoryPath, employeeLeaveRequestPath)
			if err != nil {
				log.Println("Error parsing file:", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			err = tmpl.ExecuteTemplate(w, "leaveHistory", historyData)
			if err != nil {
				log.Println("Error executing template:", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		},
	)
}

func (h *Handler) EmployeeLeaveModal() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			hxvals, err := parseRequestValues([]string{"id"}, r)
			if err != nil {
				log.Printf("Error parsing request values: %v", err)
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}

			modalData, err := application.GetLeaveRequestByIdForAdmin(hxvals[0], h.DB)
			if err != nil {
				log.Printf("Error getting admin leave modal data: %v", err)
				http.Error(w, "Not found", http.StatusNotFound)
				return
			}

			err = executePartialTemplate(employeeLeaveModalPath, "employeeLeaveModal", modalData, w)
			if err != nil {
				log.Printf("Error executing template: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		},
	)
}

func (h *Handler) RenderLeaveFormTab() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			employeeId, err := getEmployeeId(w, r)
			if err != nil {
				log.Println("Error retrieving employee_id:", err)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			employee, err := application.GetEmployeeByEmployeeId(employeeId, h.DB)
			if err != nil {
				log.Println("Error retrieving employee from db:", err)
				http.Error(w, "Status not found", http.StatusNotFound)
				return
			}

			tmpl, err := template.ParseFiles(leaveFormPath)
			if err != nil {
				log.Println("Error parsing file:", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			leaveFormDto := application.LeaveFormDto{
				EmployeeId: employee.EmployeeId,
				FirstName:  employee.FirstName,
				LastName:   employee.LastName,
				DateErr:    "",
			}

			err = tmpl.ExecuteTemplate(w, "leaveForm", leaveFormDto)
			if err != nil {
				log.Println("Error executing template:", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		},
	)
}

func (h *Handler) GetLeaveRequests() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			isAdmin, ok := r.Context().Value("is_admin").(bool)
			if !ok {
				h.Logger.Error("Error getting is_admin from context")
				http.Error(w, "Error getting value context", http.StatusNotFound)
				return
			}

			if !isAdmin {
				h.Logger.Warn("Unauthorized user")
				http.Error(w, "Unauthorized user", http.StatusUnauthorized)
				return
			}

			data, err := infrastructure.GetLeaveRequests(h.DB)
			if err != nil {
				h.Logger.Errorf("Error getting leave requests: %v", err)
				http.Error(w, "Data not found", http.StatusNotFound)
				return
			}

			responseJSON(w, data, h.Logger)
		},
	)
}

func (h *Handler) GetLeaveRequestsByEmployee() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			employeeId, ok := r.Context().Value("employee_id").(string)
			if !ok {
				h.Logger.Error("Error getting employee_id from context")
				http.Error(w, "Error getting value context", http.StatusNotFound)
				return
			}

			data, err := infrastructure.GetLeaveRequestsByEmployee(employeeId, h.DB)
			if err != nil {
				h.Logger.Errorf("Error getting leave requests for employee %s: %v", employeeId, err)
				http.Error(w, "Data not StatusNotFound", http.StatusNotFound)
				return
			}

			responseJSON(w, data, h.Logger)
		},
	)
}

func stringToDate(date string) (time.Time, error) {
	date = strings.TrimSpace(date)
	dateArr := strings.Split(date, "-")
	dateArrInt := make([]int, 3)

	for i := range dateArr {
		int, err := strconv.Atoi(dateArr[i])
		if err != nil {
			return time.Time{}, errors.New("Date value cannot be converted to integer")
		}
		dateArrInt[i] = int
	}
	month := time.Month(dateArrInt[1])
	if month < time.January || month > time.December {
		return time.Time{}, errors.New("Invalid month value")
	}

	dateObject := time.Date(dateArrInt[0], month, dateArrInt[2], 0, 0, 0, 0, time.Local)

	return dateObject, nil
}

func validateLeaveRequest(lr domain.LeaveRequest) ([]domain.ValidationResult, error) {
	results := make([]domain.ValidationResult, 0)

	from, err := stringToDate(lr.From)
	if err != nil {
		newErr := errors.New("Error converting string to date: ")
		return results, errors.Join(newErr, err)
	}
	to, err := stringToDate(lr.To)
	if err != nil {
		newErr := errors.New("Error converting string to date: ")
		return results, errors.Join(newErr, err)
	}

	if from.Before(time.Now()) {
		results = append(results, domain.ValidationResult{
			Key:     "date",
			IsValid: false,
			Msg:     "From date must be after current date",
		})
	}

	if !from.Before(to) {
		results = append(results, domain.ValidationResult{
			Key:     "date",
			IsValid: false,
			Msg:     "From date must be before To date",
		})
	}

	return results, nil
}

var postLeaveKeys = []string{"type", "from", "to", "note", "hours_per_day"}

func (h *Handler) PostLeaveRequest() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			employeeId, ok := r.Context().Value("employee_id").(string)
			if !ok {
				h.Logger.Error("Error getting employee_id from context")
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			reqVals, err := parseRequestValues(postLeaveKeys, r)
			if err != nil {
				log.Printf("Error parsing request values: %v", err)
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}

			leaveDto := application.LeaveFormDto{
				EmployeeId:  employeeId,
				Type:        reqVals[0],
				From:        reqVals[1],
				To:          reqVals[2],
				Note:        reqVals[3],
				HoursPerDay: reqVals[4],
			}

			leaveDto, err = application.PostLeaveRequest(leaveDto, h.DB)
			if err != nil {
				log.Printf("Error posting leave request: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			h.LeaveService.SendEmailNotification(&application.LeavePostNotificationHandler{})

			log.Println("email send")

			err = executePartialTemplate(leaveFormPath, "leaveForm", leaveDto, w)
			if err != nil {
				log.Println("Error executing template:", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		},
	)
}

func queryParamId(r *http.Request, w http.ResponseWriter, logger *zap.SugaredLogger) int {
	queryParams := r.URL.Query()
	idStr := queryParams.Get("id")
	if idStr == "" {
		return -1
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {

		return -1
	}

	return id
}

func (h *Handler) PutLeaveRequest() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// reading rData from request, this will be the updated rData
			var rData domain.LeaveRequest
			err := r.ParseForm()
			if err != nil {
				h.Logger.Errorf("Error parsing form: %v", err)
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}

			requestIdStr := r.FormValue("requestId")
			requestId, err := strconv.Atoi(requestIdStr)
			if err != nil {
				h.Logger.Errorf("Error converting request id to int %v:", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			rData.RequestId = requestId
			rData.Type = r.FormValue("type")
			rData.From = r.FormValue("from")
			rData.To = r.FormValue("to")
			rData.Note = r.FormValue("note")
			isApprovedStr := r.FormValue("isApproved")
			isApproved := false
			if isApprovedStr == "true" {
				isApproved = true
			}
			rData.IsApproved = isApproved

			// get leave request from db using data.RequestId
			updated, err := infrastructure.GetLeaveRequestById(rData.RequestId, h.DB)
			if err != nil {
				h.Logger.Errorf("Error getting outdated leave request %v:", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			// update old request
			updated.Type = rData.Type
			updated.From = rData.From
			updated.To = rData.To
			updated.Note = rData.Note
			updated.IsApproved = rData.IsApproved

			// validate updated leave request
			result, err := validateLeaveRequest(updated)
			fmt.Println(len(result))
			if err != nil {
				h.Logger.Errorf("Error validating leave request: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			if len(result) > 0 {
				fmt.Fprint(w, result[0].Msg)
				return
			}

			response, err := infrastructure.PutLeaveRequest(updated, h.DB)
			if err != nil {
				h.Logger.Errorf("Error updating db: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			responseJSON(w, response, h.Logger)
		},
	)
}

func (h *Handler) DeleteLeaveRequest() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			id := queryParamId(r, w, h.Logger)
			if id == -1 {
				h.Logger.Error("Error reading url parameter")
				http.Error(w, "Invalid URL parameter 'id'", http.StatusBadRequest)
				return
			}

			lr, err := infrastructure.DeleteLeaveRequest(id, h.DB)
			if err != nil {
				h.Logger.Errorf("Error deleting leave request: %v", err)
				http.Error(w, "Error updating db", http.StatusInternalServerError)
				return
			}

			responseJSON(w, lr, h.Logger)
		},
	)
}
