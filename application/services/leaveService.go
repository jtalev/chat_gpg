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

type AdminLeaveData struct {
	Pending  []domain.LeaveRequest
	Approved []domain.LeaveRequest
	Denied   []domain.LeaveRequest
}

func sortLeaveDataForAdmin(leaveRequests []domain.LeaveRequest) AdminLeaveData {
	outData := AdminLeaveData{}
	for _, lr := range leaveRequests {
		if lr.IsPending == true {
			outData.Pending = append(outData.Pending, lr)
		} else if lr.IsPending == false && lr.IsApproved == true {
			outData.Approved = append(outData.Approved, lr)
		} else if lr.IsPending == false && lr.IsApproved == false {
			outData.Denied = append(outData.Denied, lr)
		}
	}

	return outData
}

func GetLeaveRequestsForAdmin(db *sql.DB) (AdminLeaveData, error) {
	leaveRequests, err := infrastructure.GetLeaveRequests(db)
	if err != nil {
		return AdminLeaveData{}, err
	}

	adminLeaveData := sortLeaveDataForAdmin(leaveRequests)
	return adminLeaveData, nil
}
