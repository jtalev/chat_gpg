package handlers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/jtalev/chat_gpg/repository"
	"go.uber.org/zap"
)

func getLeaveData() []DashboardData {
	data := []DashboardData{
		{},
	}
	return data
}

type LeaveRequest struct {
	RequestId  int    `json:"request_id"`
	EmployeeId string `json:"employee_id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Type       string `json:"leave_type"`
	From       string `json:"from_date"`
	To         string `json:"to_date"`
	Note       string `json:"note"`
	IsApproved bool   `json:"is_approved"`
}

func PostLeaveRequest(db *sql.DB, sugar *zap.SugaredLogger) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			err := r.ParseForm()
			if err != nil {
				sugar.Errorf("Error parsing form: %v", err)
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}

			employeeId, ok := r.Context().Value("employee_id").(string)
			if !ok {
				sugar.Error("Error getting employee_id from context: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			leaveType := r.FormValue("type")
			from := r.FormValue("from")
			to := r.FormValue("to")
			note := r.FormValue("note")

			leaveRequest := LeaveRequest{
				EmployeeId: employeeId,
				Type:       leaveType,
				From:       from,
				To:         to,
				Note:       note,
			}
			leaveRequest, err = repository.PostLeaveRequest(r, leaveRequest, db, sugar)
			if err != nil {
				sugar.Errorf("Error posting leave request: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			fmt.Fprint(w, "SUBMITTED")
		},
	)
}
