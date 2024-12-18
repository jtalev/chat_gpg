package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

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

func GetLeaveRequestsByEmployee(w http.ResponseWriter, r *http.Request, db *sql.DB, sugar *zap.SugaredLogger) ([]LeaveRequest, error) {
	employeeId, ok := r.Context().Value("employee_id").(string)
	if !ok {
		fmt.Printf("employee id: %s\n", employeeId)
		return nil, errors.New("Error getting employee ID from context")
	}

	q := `
	select lr.request_id, lr.employee_id, e.first_name, e.last_name, lr.leave_type, lr.from_date,
	lr.to_date, lr.note, lr.is_approved 
	from leave_request lr 
	join employee e on lr.employee_id = e.employee_id
	where lr.employee_id = ?;
	`
	rows, err := db.Query(q, employeeId)
	if err != nil {
		return nil, errors.New("Error querying leave requests")
	}
	defer rows.Close()

	leaveRequests := []LeaveRequest{}
	for rows.Next() {
		var leaveRequest LeaveRequest
		if err := rows.Scan(
			&leaveRequest.RequestId,
			&leaveRequest.EmployeeId,
			&leaveRequest.FirstName,
			&leaveRequest.LastName,
			&leaveRequest.Type,
			&leaveRequest.From,
			&leaveRequest.To,
			&leaveRequest.Note,
			&leaveRequest.IsApproved,
		); err != nil {
			return nil, errors.New("Error scanning row")
		}
		leaveRequests = append(leaveRequests, leaveRequest)
	}

	if len(leaveRequests) == 0 {
		sugar.Infof("No leave requests found for employee %s", employeeId)
		return leaveRequests, nil
	}

	return leaveRequests, nil
}

func postLeaveRequest(r *http.Request, leaveRequest LeaveRequest, db *sql.DB, sugar *zap.SugaredLogger) (LeaveRequest, error) {
	q := `
	INSERT INTO leave_request (employee_id, leave_type, from_date, to_date, note)
	VALUES ($1, $2, $3, $4, $5);
	`
	_, err := db.Exec(q, leaveRequest.EmployeeId, leaveRequest.Type, leaveRequest.From, leaveRequest.To, leaveRequest.Note)
	if err != nil {
		return LeaveRequest{}, errors.New("Error executing db query")
	}
	return LeaveRequest{
		EmployeeId: leaveRequest.EmployeeId,
		Type:       leaveRequest.Type,
		From:       leaveRequest.From,
		To:         leaveRequest.To,
		Note:       leaveRequest.Note,
	}, nil
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
			leaveRequest, err = postLeaveRequest(r, leaveRequest, db, sugar)
			if err != nil {
				sugar.Errorf("Error posting leave request: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			fmt.Fprint(w, "SUBMITTED")
		},
	)
}
