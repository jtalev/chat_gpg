package services

import (
	"database/sql"

	"github.com/jtalev/chat_gpg/models"
	"github.com/jtalev/chat_gpg/repository"
)

func GetLeaveHistoryByEmployeeId(employeeId string, db *sql.DB) ([]models.LeaveRequest, error) {
	outLeaveHistory, err := repository.GetLeaveRequestsByEmployee(employeeId, db)
	if err != nil {
		return nil, err
	}
	return outLeaveHistory, nil
}
