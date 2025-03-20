package application

import (
	"database/sql"
	"strconv"

	"github.com/jtalev/chat_gpg/domain/models"
	"github.com/jtalev/chat_gpg/infrastructure/repository"
)

type EmployeeLeaveHistory struct {
	Pending  []domain.LeaveRequest
	Approved []domain.LeaveRequest
	Denied   []domain.LeaveRequest
}

func sortLeaveHistoryForEmployee(leaveRequests []domain.LeaveRequest) EmployeeLeaveHistory {
	outData := EmployeeLeaveHistory{}
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

func GetLeaveHistoryByEmployeeId(employeeId string, db *sql.DB) (EmployeeLeaveHistory, error) {
	empLeaveHistory, err := infrastructure.GetLeaveRequestsByEmployee(employeeId, db)
	if err != nil {
		return EmployeeLeaveHistory{}, err
	}

	employeeLeaveHistory := sortLeaveHistoryForEmployee(empLeaveHistory)

	return employeeLeaveHistory, nil
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

type AdminLeaveModalData struct {
	LeaveRequest domain.LeaveRequest
	TotalDays    int
}

func calcLeaveDaysFromRequest(lr domain.LeaveRequest) (int, error) {
	dayCounter := 1 // add one for the date the leave request starts on
	startDate, err := dateStrToDate(lr.From)
	endDate, err := dateStrToDate(lr.To)
	if err != nil {
		return -1, err
	}

	for startDate.Day() != endDate.Day() {
		startDate = startDate.AddDate(0, 0, 1)
		dayCounter += 1
	}

	return dayCounter, nil
}

func GetLeaveRequestByIdForAdmin(idStr string, db *sql.DB) (AdminLeaveModalData, error) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return AdminLeaveModalData{}, err
	}

	outLr, err := infrastructure.GetLeaveRequestById(id, db)
	if err != nil {
		return AdminLeaveModalData{}, err
	}

	totalDays, err := calcLeaveDaysFromRequest(outLr)

	outData := AdminLeaveModalData{
		outLr,
		totalDays,
	}

	return outData, nil
}

func AdminUpdateLeaveRequest(idStr, isApprovedStr string, db *sql.DB) (domain.LeaveRequest, error) {
	var isApproved bool
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return domain.LeaveRequest{}, err
	}
	if isApprovedStr == "true" {
		isApproved = true
	} else if isApprovedStr == "false" {
		isApproved = false
	} else {
		return domain.LeaveRequest{}, nil
	}

	lr, err := infrastructure.GetLeaveRequestById(id, db)
	if err != nil {
		return domain.LeaveRequest{}, nil
	}
	lr.IsApproved = isApproved
	lr.IsPending = false

	outLr, err := infrastructure.PutLeaveRequest(lr, db)
	if err != nil {
		return domain.LeaveRequest{}, nil
	}

	return outLr, nil
}

type LeaveFormDto struct {
	RequestId  string
	EmployeeId string
	FirstName  string
	LastName   string
	Type       string
	From       string
	To         string
	Note       string
	IsPending  string
	IsApproved string
	DateErr    string
	SuccessMsg string
}

func mapLeaveDtoToLeaveRequest(leaveFormDto LeaveFormDto) domain.LeaveRequest {
	leaveRequest := domain.LeaveRequest{
		EmployeeId: leaveFormDto.EmployeeId,
		Type:       leaveFormDto.Type,
		From:       leaveFormDto.From,
		To:         leaveFormDto.To,
		Note:       leaveFormDto.Note,
	}

	return leaveRequest
}

func getPendingApprovedLeave(employeeId string, db *sql.DB) ([]domain.LeaveRequest, error) {
	outRequests := []domain.LeaveRequest{}
	pastRequests, err := infrastructure.GetLeaveRequestsByEmployee(employeeId, db)
	if err != nil {
		return nil, err
	}
	for _, lr := range pastRequests {
		if lr.IsPending == false && lr.IsApproved == false {
			continue
		}
		outRequests = append(outRequests, lr)
	}
	return outRequests, nil
}

func PostLeaveRequest(leaveFormDto LeaveFormDto, db *sql.DB) (LeaveFormDto, error) {
	employee, err := infrastructure.GetEmployeeByEmployeeId(leaveFormDto.EmployeeId, db)
	if err != nil {
		return leaveFormDto, err
	}
	leaveFormDto.FirstName = employee.FirstName
	leaveFormDto.LastName = employee.LastName
	leaveRequest := mapLeaveDtoToLeaveRequest(leaveFormDto)
	pastRequests, err := getPendingApprovedLeave(leaveFormDto.EmployeeId, db)
	if err != nil {
		return leaveFormDto, err
	}

	errors := leaveRequest.Validate(pastRequests)

	if !errors.IsSuccessful {
		leaveFormDto.DateErr = errors.DateErr
		return leaveFormDto, nil
	} else {
		_, err := infrastructure.PostLeaveRequest(leaveRequest, db)
		if err != nil {
			return leaveFormDto, err
		}
		leaveFormDto.SuccessMsg = "Leave request submitted successfully."
		return leaveFormDto, nil
	}
}
