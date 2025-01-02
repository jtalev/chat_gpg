package repository

import (
	"database/sql"
	"errors"

	"github.com/jtalev/chat_gpg/models"
)

func GetLeaveRequests(db *sql.DB) ([]models.LeaveRequest, error) {
	q := `
	select lr.request_id, lr.employee_id, e.first_name, e.last_name, lr.leave_type, lr.from_date,
	lr.to_date, lr.note, lr.is_approved 
	from leave_request lr 
	join employee e on lr.employee_id = e.employee_id
	order by e.employee_id asc;
	`
	rows, err := db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := []models.LeaveRequest{}
	for rows.Next() {
		leaveRequest := models.LeaveRequest{}
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
			return nil, err
		}
		data = append(data, leaveRequest)
	}
	return data, nil
}

func GetLeaveRequestById(requestId int, db *sql.DB) (models.LeaveRequest, error) {
	q := `
	select lr.request_id, lr.employee_id, e.first_name, e.last_name, lr.leave_type, lr.from_date,
	lr.to_date, lr.note, lr.is_approved 
	from leave_request lr 
	join employee e on lr.employee_id = e.employee_id
	where lr.request_id = ?;
	`
	rows, err := db.Query(q, requestId)
	if err != nil {
		return models.LeaveRequest{}, err
	}
	defer rows.Close()

	var data models.LeaveRequest
	if rows.Next() {
		if err := rows.Scan(
			&data.RequestId,
			&data.EmployeeId,
			&data.FirstName,
			&data.LastName,
			&data.Type,
			&data.From,
			&data.To,
			&data.Note,
			&data.IsApproved,
		); err != nil {
			return data, err
		}
	} else {
		return data, errors.New("No leave request with provided requestId")
	}

	return data, nil
}

func GetLeaveRequestsByEmployee(employeeId string, db *sql.DB) ([]models.LeaveRequest, error) {
	q := `
	select lr.request_id, lr.employee_id, e.first_name, e.last_name, lr.leave_type, lr.from_date,
	lr.to_date, lr.note, lr.is_approved 
	from leave_request lr 
	join employee e on lr.employee_id = e.employee_id
	where lr.employee_id = ?;
	`
	rows, err := db.Query(q, employeeId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	leaveRequests := []models.LeaveRequest{}
	for rows.Next() {
		var leaveRequest models.LeaveRequest
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
			return nil, err
		}
		leaveRequests = append(leaveRequests, leaveRequest)
	}

	return leaveRequests, nil
}

func PostLeaveRequest(leaveRequest models.LeaveRequest, db *sql.DB) (models.LeaveRequest, error) {
	q := `
	INSERT INTO leave_request (employee_id, leave_type, from_date, to_date, note)
	VALUES ($1, $2, $3, $4, $5);
	`
	_, err := db.Exec(q, leaveRequest.EmployeeId, leaveRequest.Type, leaveRequest.From, leaveRequest.To, leaveRequest.Note)
	if err != nil {
		return models.LeaveRequest{}, err
	}
	return models.LeaveRequest{
		EmployeeId: leaveRequest.EmployeeId,
		Type:       leaveRequest.Type,
		From:       leaveRequest.From,
		To:         leaveRequest.To,
		Note:       leaveRequest.Note,
	}, nil
}

func PutLeaveRequest(leaveRequest models.LeaveRequest, db *sql.DB) (models.LeaveRequest, error) {
	q := `
	update leave_request
	set leave_type = $1, from_date = $2, to_date = $3, note = $4, is_approved = $5
	where request_id = $6;
	`

	_, err := db.Exec(
		q,
		leaveRequest.Type,
		leaveRequest.From,
		leaveRequest.To,
		leaveRequest.Note,
		leaveRequest.IsApproved,
		leaveRequest.RequestId,
	)
	if err != nil {
		return models.LeaveRequest{}, err
	}
	return leaveRequest, nil
}

func DeleteLeaveRequest(requestId int, db *sql.DB) (models.LeaveRequest, error) {
	q := `
	delete from leave_request where request_id = ?;
	`
	_, err := db.Exec(q, requestId)
	if err != nil {
		return models.LeaveRequest{}, err
	}

	lr, err := GetLeaveRequestById(requestId, db)
	if err == nil {
		return lr, errors.New("Leave request still exists")
	}

	return lr, nil
}
