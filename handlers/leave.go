package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"go.uber.org/zap"
)

type LeaveRequest struct {
	EmployeeId int    `json:"employee_id"`
	RequestId  int    `json:"request_id"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	From       string `json:"from"`
	To         string `json:"to"`
	Note       string `json:"note"`
	IsApproved bool   `json:"is_approved"`
}

var tempData = []LeaveRequest{}

func ServeLeaveView(sugar *zap.SugaredLogger) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			data := GetLeaveRequestById(sugar)
			component := "leave"
			title := "Leave - GPG"
			renderTemplate(w, component, title, data)
		},
	)
}

func GetLeaveRequests(db *sql.DB, sugar *zap.SugaredLogger) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			sugar.Info("processing leave GET request")
			responseJson(w, tempData, sugar)
		},
	)
}

func GetLeaveRequestById(sugar *zap.SugaredLogger) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			sugar.Info("processing leave GET request")

			if err := r.ParseForm(); err != nil {
				sugar.Error(err)
				http.Error(w, "error parsing form", http.StatusBadRequest)
				return
			}

			formRequestId := r.FormValue("requestId")
			if formRequestId == "" {
				http.Error(w, "form request id empty string", http.StatusBadRequest)
				sugar.Error("form request id empty string")
				return
			}
			requestId, err := strconv.Atoi(formRequestId)
			if err != nil {
				http.Error(w, "invalid request id", http.StatusNotFound)
				sugar.Errorf("invalid request id: %v", err)
				return
			}

			found := false
			var requestIndex int
			for i, _ := range tempData {
				if tempData[i].RequestId == requestId {
					requestIndex = i
					found = true
				}
			}

			if !found {
				sugar.Errorf("leave request not found: %v", err)
				http.Error(w, "leave request not found", http.StatusNotFound)
				return
			}
			responseJson(w, tempData[requestIndex], sugar)
		},
	)
}

func PostLeaveRequest(sugar *zap.SugaredLogger) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			sugar.Info("processing leave POST request")

			if err := r.ParseForm(); err != nil {
				sugar.Error(err)
				http.Error(w, "error parsing form", http.StatusBadRequest)
				return
			}

			name := r.FormValue("name")
			from := r.FormValue("from")
			to := r.FormValue("to")
			note := r.FormValue("note")
			leaveType := r.FormValue("type")

			leaveRequest := LeaveRequest{
				EmployeeId: 12345678,
				RequestId:  87654321,
				Name:       name,
				Type:       leaveType,
				From:       from,
				To:         to,
				Note:       note,
				IsApproved: false,
			}

			// TODO: send leave request to db
			tempData = append(tempData, leaveRequest)
			fmt.Println(len(tempData))

			sugar.Infof("successfully submitted leave request: %v", tempData[len(tempData)-1])
			fmt.Fprintf(w, "SUBMITTED")
		},
	)
}

func UpdateLeaveRequest(sugar *zap.SugaredLogger) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			sugar.Info("processing leave PUT request")

			if err := r.ParseForm(); err != nil {
				sugar.Error(err)
				http.Error(w, "error parsing form", http.StatusBadRequest)
				return
			}

			type UpdatedRequestData struct {
				Request      LeaveRequest
				RequestIndex int
			}
			var updatedRequestData UpdatedRequestData

			formRequestId := r.FormValue("requestId")
			if formRequestId == "" {
				http.Error(w, "form request id empty string", http.StatusBadRequest)
				sugar.Error("form request id empty string")
				return
			}
			requestId, err := strconv.Atoi(formRequestId)
			if err != nil {
				http.Error(w, "invalid request id", http.StatusNotFound)
				sugar.Errorf("invalid request id: %v", err)
				return
			}
			from := r.FormValue("from")
			to := r.FormValue("to")
			note := r.FormValue("note")
			leaveType := r.FormValue("type")

			found := false
			for i, _ := range tempData {
				if tempData[i].RequestId == requestId {
					updatedRequestData.Request = tempData[i]
					updatedRequestData.RequestIndex = i
					found = true
				}
			}

			formRequest := updatedRequestData.Request

			leaveRequest := LeaveRequest{
				EmployeeId: formRequest.EmployeeId,
				RequestId:  formRequest.RequestId,
				Name:       formRequest.Name,
				Type:       leaveType,
				From:       from,
				To:         to,
				Note:       note,
				IsApproved: formRequest.IsApproved,
			}

			if !found {
				http.Error(w, "leave request not found", http.StatusBadRequest)
				return
			}

			tempData[updatedRequestData.RequestIndex] = leaveRequest
			fmt.Println(tempData[updatedRequestData.RequestIndex])

			w.WriteHeader(http.StatusOK)
			responseJson(w, leaveRequest, sugar)
		},
	)
}

func DeleteLeaveRequest(sugar *zap.SugaredLogger) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			sugar.Info("processing leave DELETE request")

			if err := r.ParseForm(); err != nil {
				sugar.Error(err)
				http.Error(w, "error parsing form", http.StatusBadRequest)
			}

			formRequestId := r.FormValue("requestId")
			if formRequestId == "" {
				http.Error(w, "form request id empty string", http.StatusBadRequest)
				sugar.Error("form request id empty string")
				return
			}
			requestId, err := strconv.Atoi(formRequestId)
			if err != nil {
				http.Error(w, "invalid request id", http.StatusNotFound)
				sugar.Errorf("invalid request id: %v", err)
				return
			}

			found := false

			var requestIndex int
			for i, _ := range tempData {
				if tempData[i].RequestId == requestId {
					tempData = append(tempData[:i], tempData[i+1:]...)
					requestIndex = i
					found = true
					break
				}
			}

			if !found {
				http.Error(w, "leave request not found", http.StatusNotFound)
				sugar.Error("leave request not found")
				return
			}

			w.WriteHeader(http.StatusOK)
			responseJson(w, tempData[requestIndex], sugar)
		},
	)
}
