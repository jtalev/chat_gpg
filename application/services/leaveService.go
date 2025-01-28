package application

import (
	"database/sql"

	"github.com/jtalev/chat_gpg/domain/models"
	"github.com/jtalev/chat_gpg/infrastructure/repository"
)

func GetLeaveHistoryByEmployeeId(employeeId string, db *sql.DB) ([]domain.LeaveRequest, error) {
	outLeaveHistory, err := infrastructure.GetLeaveRequestsByEmployee(employeeId, db)
	if err != nil {
		return nil, err
	}
	return outLeaveHistory, nil
}
