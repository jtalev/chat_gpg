package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/jtalev/chat_gpg/models"
	"github.com/jtalev/chat_gpg/repository"
	"go.uber.org/zap"
)

func getLeaveData() []DashboardData {
	data := []DashboardData{
		{},
	}
	return data
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

			data, err := repository.GetLeaveRequests(isAdmin, h.DB)
			if err != nil {
				h.Logger.Errorf("Error getting leave requests: %v", err)
				http.Error(w, "Data not found", http.StatusNotFound)
				return
			}

			responseJson(w, data, h.Logger)
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

			data, err := repository.GetLeaveRequestsByEmployee(employeeId, h.DB)
			if err != nil {
				h.Logger.Errorf("Error getting leave requests for employee %s: %v", employeeId, err)
				http.Error(w, "Data not StatusNotFound", http.StatusNotFound)
			}

			responseJson(w, data, h.Logger)
		},
	)
}

func (h *Handler) PostLeaveRequest() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			err := r.ParseForm()
			if err != nil {
				h.Logger.Errorf("Error parsing form: %v", err)
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}

			employeeId, ok := r.Context().Value("employee_id").(string)
			if !ok {
				h.Logger.Error("Error getting employee_id from context")
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			leaveType := r.FormValue("type")
			from := r.FormValue("from")
			to := r.FormValue("to")
			note := r.FormValue("note")

			leaveRequest := models.LeaveRequest{
				EmployeeId: employeeId,
				Type:       leaveType,
				From:       from,
				To:         to,
				Note:       note,
			}
			leaveRequest, err = repository.PostLeaveRequest(leaveRequest, h.DB)
			if err != nil {
				h.Logger.Errorf("Error posting leave request: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			fmt.Fprint(w, "SUBMITTED")
		},
	)
}

func queryParamId(r *http.Request, w http.ResponseWriter, logger *zap.SugaredLogger) int {
	queryParams := r.URL.Query()
	logger.Info(queryParams)
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
			var rData models.LeaveRequest
			if err := json.NewDecoder(r.Body).Decode(&rData); err != nil {
				http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
				return
			}
			h.Logger.Info(rData)

			// get url parameter and add to data.RequestId
			id := queryParamId(r, w, h.Logger)
			rData.RequestId = id
			h.Logger.Info(rData.RequestId)
			if id == -1 {
				h.Logger.Error("Error reading url parameter")
				http.Error(w, "Invalid URL parameter 'id'", http.StatusBadRequest)
				return
			}

			// get leave request from db using data.RequestId
			updated, err := repository.GetLeaveRequestById(rData.RequestId, h.DB)
			if err != nil {
				h.Logger.Errorf("Error getting outdated leave request %v:", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			h.Logger.Info(updated)

			// update old request
			updated.Type = rData.Type
			updated.From = rData.From
			updated.To = rData.To
			updated.Note = rData.Note
			updated.IsApproved = rData.IsApproved

			h.Logger.Info(updated)
			response, err := repository.PutLeaveRequest(updated, h.DB)
			if err != nil {
				h.Logger.Errorf("Error updating db: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			responseJson(w, response, h.Logger)
		},

		// update db with new updated leave request
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

			lr, err := repository.DeleteLeaveRequest(id, h.DB)
			if err != nil {
				h.Logger.Errorf("Error deleting leave request: %v", err)
				http.Error(w, "Error updating db", http.StatusInternalServerError)
				return
			}

			responseJson(w, lr, h.Logger)
		},
	)
}
